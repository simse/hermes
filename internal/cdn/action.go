package cdn

import (
	"github.com/simse/hermes/internal/action"
)

func init() {
	// Register all actions
	createOAIAction := action.Action{
		Name:            "oai:create",
		ProgressMessage: "Creating Origin Access Identity...",
		FinishedMessage: "Created Origin Access Identity.",
		Handler:         CreateOAIAction,
	}

	action.AddAction(createOAIAction)
}

// CreateOAIAction represents a function that creates an OAI
func CreateOAIAction(input action.Input) action.Output {
	_, err := CreateOAI(
		input.Payload["oai:comment"].(string),
	)

	if err != nil {
		return action.Output{
			Status: action.ERROR,
		}
	}

	return action.Output{
		Status: action.OK,
	}
}

// CreateDistributionAction represents an action that creates a CF distribution
func CreateDistributionAction(input action.Input) action.Output {
	return action.Output{
		Status: action.OK,
	}
}
