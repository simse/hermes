package edge

import (
	"time"

	"github.com/simse/hermes/internal/action"
)

func init() {
	action.AddAction(action.Action{
		Handler:         CreateExecutionRoleAction,
		Name:            "execution_role:create",
		ProgressMessage: "Creating lambda@edge execution role...",
		FinishedMessage: "Created lambda@edge execution role",
	})

	action.AddAction(action.Action{
		Handler:         CreateLambdaFunctionAction,
		Name:            "edge_function:create",
		ProgressMessage: "Creating lambda function...",
		FinishedMessage: "Created lambda function",
	})
}

// CreateExecutionRoleAction creates an lambda execution role
func CreateExecutionRoleAction(input action.Input) action.Output {
	role, roleCreateErr := CreateExecutionRole(
		input.Payload["execution_role:name"].(string),
	)

	input.Environment["iam:role_arn"] = role

	if roleCreateErr != nil {
		return action.Output{
			Status:      action.ERROR,
			Environment: input.Environment,
		}
	}

	// Due an apparent bug in the lambda api, it's neccessary to wait a bit before moving on
	time.Sleep(time.Second * 8)

	return action.Output{
		Status:      action.OK,
		Environment: input.Environment,
	}
}

// CreateLambdaFunctionAction creates and optionally publishes a lambda function
func CreateLambdaFunctionAction(input action.Input) action.Output {
	lambdaError := CreateLambdaFunction(
		input.Payload["lambda:name"].(string),
		input.Environment["iam:role_arn"].(string),
		input.Payload["lambda:runtime"].(string),
		input.Payload["lambda:func_path"].(string),
	)

	if lambdaError != nil {
		return action.Output{
			Status:      action.ERROR,
			Environment: input.Environment,
		}
	}

	config, publishError := PublishLambdaFunction(input.Payload["lambda:name"].(string))
	if publishError != nil {
		return action.Output{
			Status:      action.ERROR,
			Environment: input.Environment,
		}
	}

	/*if _, ok := input.Environment["lambda:edge_handlers"]; !ok {
		input.Environment["lambda:edge_handlers"] = map[string]string{}
	}*/

	input.Environment["lambda:edge_handlers:"+input.Payload["lambda:event"].(string)] = *config.FunctionArn

	return action.Output{
		Status:      action.OK,
		Environment: input.Environment,
	}
}
