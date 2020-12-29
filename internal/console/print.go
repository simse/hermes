package console

import (
	"fmt"
)

// ShowUsing formats an action sheet action as using
func ShowUsing(service, name string) {
	fmt.Print("> [" + service + "] ")
	WhiteBold.Print(name)
	fmt.Print("\n")
}

// ShowDelete formats an action sheet delete
func ShowDelete(service, name string) {
	Red.Print("-")
	fmt.Print(" [" + service + "] ")
	WhiteBold.Print(name)
	fmt.Print("\n")
}

// ShowCreate formats an action sheet create
func ShowCreate(service, name string) {
	Green.Print("+")
	fmt.Print(" [" + service + "] ")
	WhiteBold.Print(name)
	fmt.Print("\n")
}

// ShowAttach shows an action sheet attach
func ShowAttach(firstService, firstName, secondService, secondName string) {
	Cyan.Print("|")
	fmt.Print(" [" + firstService + "] ")
	WhiteBold.Print(firstName)
	Cyan.Print(" ->")
	fmt.Print(" [" + secondService + "] ")
	WhiteBold.Print(secondName)
	fmt.Print("\n")
}

// ShowLegend explains the symbols
func ShowLegend() {
	fmt.Print("> using, ")
	Green.Print("+")
	fmt.Print(" create, ")
	Red.Print("-")
	fmt.Print(" delete, ")
	Cyan.Print("->")
	fmt.Print(" attach resource")
	fmt.Print("\n\n")
}
