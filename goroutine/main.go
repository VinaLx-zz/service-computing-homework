package main

import (
	"fmt"
	"scenario1"
	"scenario2"

	"github.com/ogier/pflag"
)

var choice *string

func init() {
	choice = pflag.StringP("scenario", "s", "",
		"choose a scenario, 1 for the synchronized scenario,"+
			" 2 for the asynchronized one")
	pflag.Parse()
}

func main() {
	switch *choice {
	case "1":
		scenario1.Start()
	case "2":
		scenario2.Start()
	default:
		fmt.Println("choose a scenario")
	}
}
