package bucket

import (
	"strings"

	"github.com/simse/hermes/internal/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 stores the current connection to S3
var S3 *s3.S3

// InitS3 intialises an S3 service
func InitS3() {
	svc := s3.New(session.Session)
	S3 = svc
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
	return nil
}

// AddOAIPermissions adds a CloudFront Origin Access Identity read permission to the bucket
func AddOAIPermissions(bucket, canonicalUser string) error {
	return nil
}
