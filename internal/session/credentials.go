package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session contains the current AWS connection
var Session *session.Session

// SecondarySession may contain a session in a different region
var SecondarySession *session.Session

// InitSession connects to AWS
func InitSession() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		panic(err)
	}

	Session = sess
}

// InitSecondarySession connects to AWS
func InitSecondarySession(region string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		panic(err)
	}

	Session = sess
}
