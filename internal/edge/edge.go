package edge

import (
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/simse/hermes/internal/session"
)

// Lambda stores the current connection to AWS lambda
var Lambda *lambda.Lambda

// InitLambda creates a lambda service
func InitLambda() {
	Lambda = lambda.New(session.Session)
}

// CreateLambdaFunction creates a lambda function
func CreateLambdaFunction(name, region, runtime, contents string) error {
	return nil
}

// CreateExecutionRole creates a role so the function works with lambda@edge
func CreateExecutionRole(name string) (string, error) {
	return "", nil
}
