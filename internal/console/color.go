package console

import (
	"github.com/fatih/color"
)

// WhiteBold is a bold white
var WhiteBold = color.New(color.FgWhite).Add(color.Bold)

// WhiteUnderline is a bold white with an underline
var WhiteUnderline = color.New(color.FgWhite).Add(color.Bold).Add(color.Underline)

// Green is green
var Green = color.New(color.FgGreen)

// Red is red
var Red = color.New(color.FgRed)

// Cyan is cyan
var Cyan = color.New(color.FgCyan)
