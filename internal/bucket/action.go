package bucket

import (
	"github.com/simse/hermes/internal/action"
)

func init() {
	// Register all actions
	createBucketAction := action.Action{
		Name:            "bucket:create",
		ProgressMessage: "Creating bucket...",
		FinishedMessage: "Created bucket.",
		Handler:         CreateBucketAction,
	}

	addOAIAction := action.Action{
		Name:            "bucket:add_oai",
		ProgressMessage: "Adding Origin Access Identity to bucket...",
		FinishedMessage: "Added Origin Access Identity to bucket.",
		Handler:         AddOriginAccessIdentityAction,
	}

	action.AddAction(createBucketAction)
	action.AddAction(addOAIAction)
}

// CreateBucketAction represents a function that creates a S3 bucket
func CreateBucketAction(input action.Input) action.Output {
	createBucketErr := Create(
		input.Payload["bucket:name"].(string),
		input.Payload["bucket:region"].(string),
		input.Payload["bucket:public"].(bool),
	)

	input.Environment["bucket:name"] = input.Payload["bucket:name"].(string)

	if createBucketErr != nil {
		return action.Output{
			Status:      action.ERROR,
			Environment: input.Environment,
		}
	}

	return action.Output{
		Status:      action.OK,
		Environment: input.Environment,
	}
}

// AddOriginAccessIdentityAction adds an OAI policy to a bucket
func AddOriginAccessIdentityAction(input action.Input) action.Output {
	addOAIErr := AddOAIPermissions(
		input.Payload["bucket:name"].(string),
		input.Payload["oai:canonical_user"].(string),
	)

	if addOAIErr != nil {
		return action.Output{
			Status: action.ERROR,
		}
	}

	return action.Output{
		Status: action.OK,
	}
}
