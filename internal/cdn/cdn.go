package cdn

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
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
