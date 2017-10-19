package cli

import (
	"os"
	"fmt"
)

func printError(error string) {
	fmt.Fprint(os.Stderr, error)
	os.Exit(1)
}

func checkEmpty(key, value string) {
	if value == "" {
		printError(key + " can't be empty!\n")
	}
}
