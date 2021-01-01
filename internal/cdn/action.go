package cdn

import (
	"strings"

	"github.com/simse/hermes/internal/action"
)

func init() {
	// Register all actions
	createOAIAction := action.Action{
		Name:            "oai:create",
		ProgressMessage: "Creating Origin Access Identity...",
		FinishedMessage: "Created Origin Access Identity",
		Handler:         CreateOAIAction,
	}

	createDistribution := action.Action{
		Name:            "cdn:create",
		ProgressMessage: "Creating CloudFront distribution...",
		FinishedMessage: "Created CloudFront distribution",
		Handler:         CreateDistributionAction,
	}

	action.AddAction(createOAIAction)
	action.AddAction(createDistribution)
}

// CreateOAIAction represents a function that creates an OAI
func CreateOAIAction(input action.Input) action.Output {
	identity, err := CreateOAI(
		input.Payload["oai:comment"].(string),
	)

	input.Environment["oai:canonical_user"] = *identity.S3CanonicalUserId
	input.Environment["oai:id"] = *identity.Id

	if err != nil {
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

// CreateDistributionAction represents an action that creates a CF distribution
func CreateDistributionAction(input action.Input) action.Output {
	lambdaMappings := map[string]string{}
	// Get lambda function mappings
	for key, value := range input.Environment {
		if strings.HasPrefix(key, "lambda:edge_handlers:") {
			event := strings.TrimPrefix(key, "lambda:edge_handlers:")
			lambdaMappings[event] = value.(string)
		}
	}

	distribution, distributionErr := CreateDistribution(
		input.Payload["cdn:comment"].(string),
		input.Environment["bucket:name"].(string),
		input.Environment["oai:id"].(string),
		"PriceClass_All",
		[]string{input.Environment["domain"].(string)},
		lambdaMappings,
	)

	if distributionErr != nil {
		return action.Output{
			Status:      action.ERROR,
			Environment: input.Environment,
		}
	}

	input.Environment["cdn:id"] = *distribution.Id
	input.Environment["cdn:arm"] = *distribution.ARN
	input.Environment["cdn:domain"] = *distribution.DomainName

	return action.Output{
		Status:      action.OK,
		Environment: input.Environment,
	}
}
