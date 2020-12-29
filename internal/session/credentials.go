package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session contains the current AWS connection
var Session *session.Session

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
