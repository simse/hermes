package action

// Input represents an action like creating a bucket
type Input struct {
	ID      string
	Handler string
	Payload map[string]string
}

// Output represents an action output
type Output struct {
	Status string
}

// Action defines an action handler
type Action struct {
	Handler         func(Input) (Output, error)
	ProgressMessage string
	FinishedMessage string
}
