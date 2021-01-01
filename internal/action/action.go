package action

// Input represents an action like creating a bucket
type Input struct {
	ID          string
	Handler     string
	Payload     map[string]interface{}
	Environment map[string]interface{}
}

// Output represents an action output
type Output struct {
	Status      string
	Payload     map[string]string
	Environment map[string]interface{}
}

// Action defines an action handler
type Action struct {
	Handler         func(Input) Output
	ProgressMessage string
	FinishedMessage string
	ErrorMessage    string
	Name            string
}

// Constants
var (
	OK    = "OK"
	ERROR = "ERROR"
)

// Actions stores all available actions
var Actions = map[string]Action{}

// AddAction adds an action handle to the registry
func AddAction(action Action) {
	Actions[action.Name] = action
}

// Run runs an action given input
func Run(input Input) Output {
	return Actions[input.Handler].Handler(input)
}
