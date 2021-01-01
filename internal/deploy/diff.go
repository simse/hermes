package deploy

// Diff describes how to get from one folder state to another
type Diff struct {
	Add    []File
	Rename []Rename
	Delete []Delete
}

// Rename represents a rename action
type Rename struct {
	From string
	To   string
}

// Delete represents a delete action
type Delete struct {
	Key string
}
