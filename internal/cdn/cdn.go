package cdn

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/simse/hermes/internal/certificate"
	"github.com/simse/hermes/internal/session"
)

// CloudFront contains a connection to AWS CloudFront
var CloudFront *cloudfront.CloudFront

// InitCloudFront creates and stores a connection to CloudFront
func InitCloudFront() {
	CloudFront = cloudfront.New(session.Session)
}

// CheckConflictCNAME checks if there's a CloudFront distribution alread using the CNAME
func CheckConflictCNAME(domain string) bool {
	// List all distributions
	listDistributionsInput := cloudfront.ListDistributionsInput{}
	listDistributionsOutput, err := CloudFront.ListDistributions(&listDistributionsInput)
	if err != nil {
		panic(err)
	}

	// Check every distribution individually
	for _, dist := range listDistributionsOutput.DistributionList.Items {
		// Check every alias
		for _, alias := range dist.Aliases.Items {
			if *alias == domain {
				return true
			}
		}
	}

	return false
}

// CreateOAI will create an Origin Access Identity in CloudFront
func CreateOAI(comment string) (*cloudfront.OriginAccessIdentity, error) {
	OAIRequestInput := cloudfront.CreateCloudFrontOriginAccessIdentityInput{
		CloudFrontOriginAccessIdentityConfig: &cloudfront.OriginAccessIdentityConfig{
			CallerReference: aws.String(time.Now().Format("2006-01-02 15:04:05.000000000")), // Long date
			Comment:         aws.String(comment),
		},
	}

	OAIRequestOutput, err := CloudFront.CreateCloudFrontOriginAccessIdentity(&OAIRequestInput)
	if err != nil {
		panic(err)
	}

	return OAIRequestOutput.CloudFrontOriginAccessIdentity, nil

	//return "", nil
}

// CreateDistribution will create a CloudFront distribution
func CreateDistribution(comment, bucket, OAI, priceClass string, aliasDomains []string, lambdas map[string]string) error {
	s3OriginConfig := cloudfront.S3OriginConfig{
		OriginAccessIdentity: aws.String("origin-access-identity/cloudfront/" + OAI),
	}

	originID := aws.String("S3-" + bucket)

	origin := cloudfront.Origin{
		DomainName:     aws.String(bucket + ".s3.amazonaws.com"),
		Id:             originID,
		S3OriginConfig: &s3OriginConfig,
	}

	origins := cloudfront.Origins{
		Items:    []*cloudfront.Origin{&origin},
		Quantity: aws.Int64(1),
	}

	lambdaFunctionAssociations := cloudfront.LambdaFunctionAssociations{
		Quantity: aws.Int64(0),
	}
	for event, handler := range lambdas {
		lambdaFunctionAssociation := cloudfront.LambdaFunctionAssociation{
			EventType:         aws.String(event),
			LambdaFunctionARN: aws.String(handler),
		}

		lambdaFunctionAssociations.Items = append(lambdaFunctionAssociations.Items, &lambdaFunctionAssociation)
		lambdaFunctionAssociations.Quantity = aws.Int64(*lambdaFunctionAssociations.Quantity + 1)
	}

	defaultCacheBehavior := cloudfront.DefaultCacheBehavior{
		TargetOriginId:             originID,
		ViewerProtocolPolicy:       aws.String("redirect-to-https"),
		CachePolicyId:              aws.String("658327ea-f89d-4fab-a63d-7e88639e58f6"),
		Compress:                   aws.Bool(true),
		LambdaFunctionAssociations: &lambdaFunctionAssociations,
	}

	aliases := cloudfront.Aliases{
		Quantity: aws.Int64(0),
	}
	for _, alias := range aliasDomains {
		aliases.Items = append(aliases.Items, aws.String(alias))
		aliases.Quantity = aws.Int64(*aliases.Quantity + 1)
	}

	// Find certificate
	certificate, err := certificate.Get(aliasDomains[0])
	if err != nil {
		panic(err)
	}

	viewerCertificate := cloudfront.ViewerCertificate{
		ACMCertificateArn:      certificate.CertificateArn,
		MinimumProtocolVersion: aws.String("TLSv1"),
		SSLSupportMethod:       aws.String("sni-only"),
	}

	distributionConfig := cloudfront.DistributionConfig{
		Comment:              aws.String(comment),
		Enabled:              aws.Bool(true),
		PriceClass:           aws.String(priceClass),
		Origins:              &origins,
		DefaultRootObject:    aws.String("index.html"),
		DefaultCacheBehavior: &defaultCacheBehavior,
		CallerReference:      aws.String(time.Now().Format("2006-01-02 15:04:05.000000000")),
		Aliases:              &aliases,
		ViewerCertificate:    &viewerCertificate,
	}

	tag := cloudfront.Tag{
		Key:   aws.String("X-Created-By"),
		Value: aws.String("hermes"),
	}

	tags := cloudfront.Tags{
		Items: []*cloudfront.Tag{&tag},
	}

	distributionConfigWithTags := cloudfront.DistributionConfigWithTags{
		DistributionConfig: &distributionConfig,
		Tags:               &tags,
	}

	createDistributionInput := cloudfront.CreateDistributionWithTagsInput{
		DistributionConfigWithTags: &distributionConfigWithTags,
	}

	createDistributionOutput, err := CloudFront.CreateDistributionWithTags(&createDistributionInput)
	if err != nil {
		panic(err)
	}

	fmt.Println(createDistributionOutput)

	return nil
}
