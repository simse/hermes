package action

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/simse/hermes/internal/console"
)

// List represents a list of actions
type List struct {
	Actions     []Input
	Environment map[string]interface{}
}

// Add adds an item to the action list
func (l *List) Add(name string, payload map[string]interface{}) error {
	l.Actions = append(l.Actions, Input{
		Handler: name,
		Payload: payload,
		ID:      randString(64),
	})

	return nil
}

// RunAll runs every action given a slice of inputs
func (l *List) RunAll() []Output {
	var outputs []Output

	for _, input := range l.Actions {
		action := Actions[input.Handler]
		input.Environment = l.Environment

		actionSpinner := wow.New(os.Stdout, spin.Get(spin.Dots), " "+action.ProgressMessage)
		actionSpinner.Start()

		actionOutput := Run(input)

		outputs = append(outputs, actionOutput)
		l.Environment = actionOutput.Environment

		if actionOutput.Status == OK {
			actionSpinner.PersistWith(console.Check, " "+action.FinishedMessage+"\n")
		} else {
			actionSpinner.PersistWith(console.Cross, " "+action.ErrorMessage+"\n")
		}
	}

	return outputs
}

func randString(n int) string {
	rand.Seed(time.Now().Unix())
	var output strings.Builder

	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	for i := 0; i < n; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}
