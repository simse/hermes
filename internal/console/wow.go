package console

import (
	"github.com/fatih/color"
	"github.com/gernest/wow/spin"
)

var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

// Check contains a checkmark for use with wow spinner
var Check = spin.Spinner{Frames: []string{green("✓")}}

// Cross contains an error cross for use with wow spinner
var Cross = spin.Spinner{Frames: []string{red("✗")}}
