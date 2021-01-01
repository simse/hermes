package bucket

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/markbates/pkger"
	"github.com/simse/hermes/internal/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 stores the current connection to S3
var S3 *s3.S3

// InitS3 intialises an S3 service
func InitS3() {
	S3 = s3.New(session.Session)
}

// Exists checks if bucket name is available
func Exists(name string) (bool, error) {
	// Get some info about the bucket
	headBucketInput := s3.HeadBucketInput{
		Bucket: aws.String(name),
	}

	_, headBucketError := S3.HeadBucket(&headBucketInput)

	// If there's an error and its 404, bucket does not exist
	if headBucketError != nil {
		// fmt.Println(headBucketError.Error())

		if strings.HasPrefix(headBucketError.Error(), "NotFound") {
			return false, nil
		}

		return true, ErrBucketExistsForeign
	}

	// Response code was 200 (our bucket) or 403 (someone else's bucket)
	return true, ErrBucketExists
}

// Create creates a bucket in the specified region
func Create(name, region string, public bool) error {
	// TODO: Allow public bucket access
	if public {
		panic("Unimplemented feature")
	}

	var createBucketInput s3.CreateBucketInput

	// If region is not default, explicitly define it
	if region != "us-east-1" {
		createBucketInput = s3.CreateBucketInput{
			Bucket: aws.String(name),
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{
				LocationConstraint: aws.String(region),
			},
		}
	} else {
		createBucketInput = s3.CreateBucketInput{
			Bucket: aws.String(name),
		}
	}

	// Perform create action
	_, err := S3.CreateBucket(&createBucketInput)
	if err != nil {
		panic(err)
	}

	// Wait for bucket to be created
	getBucketInput := s3.HeadBucketInput{
		Bucket: aws.String(name),
	}

	S3.WaitUntilBucketExists(&getBucketInput)

	// Set bucket X-Created-By tag
	tagSet := []*s3.Tag{
		&s3.Tag{
			Key:   aws.String("X-Created-By"),
			Value: aws.String("hermes"),
		},
	}

	tagging := s3.Tagging{
		TagSet: tagSet,
	}

	putBucketTaggingInput := s3.PutBucketTaggingInput{
		Bucket:  aws.String(name),
		Tagging: &tagging,
	}

	_, taggingErr := S3.PutBucketTagging(&putBucketTaggingInput)
	if taggingErr != nil {
		panic(taggingErr)
	}

	return nil
}

// AddOAIPermissions adds a CloudFront Origin Access Identity read permission to the bucket
func AddOAIPermissions(bucket, canonicalUser string) error {
	policyFile, openError := pkger.Open("/assets/s3/bucketPolicy.json")
	if openError != nil {
		panic(openError)
	}

	policyTemplate, _ := ioutil.ReadAll(policyFile)

	policy := bytes.ReplaceAll(policyTemplate, []byte("{{ CANONCICAL_USER }}"), []byte(canonicalUser))
	policy = bytes.Replace(policy, []byte("{{ BUCKET }}"), []byte(bucket), 1)

	putBucketPolicyInput := s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(string(policy)),
	}

	_, err := S3.PutBucketPolicy(&putBucketPolicyInput)
	if err != nil {
		panic(err)
	}

	//fmt.Println(putBucketPolicyOutput)

	return nil
}

// GenerateAlternativeNames will generate a list of available bucket names
func GenerateAlternativeNames(inspiration string) []string {
	potentialNames := []string{
		inspiration + "-2",
		inspiration + "-3",
		inspiration + "-4",
		inspiration + "-hermes",
		inspiration + "-hermes-2",
		inspiration + "-hermes-3",
	}
	availableNames := []string{}

	for _, potentialName := range potentialNames {
		bucketExists, _ := Exists(potentialName)

		if !bucketExists {
			availableNames = append(availableNames, potentialName)
		}
	}

	return availableNames
}
