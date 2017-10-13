package main

import (
	"args"
	"fmt"
	"os"
)

func errorExit(reason string) {
	fmt.Fprintln(os.Stderr, reason)
	os.Exit(1)
}

func selpg(a *args.Args) {
	for _, reader := range a.Sources {
		n := 0
		for page := range a.Pager(reader) {
			yes, last := a.Filter(n, page)
			if yes {
				_, err := a.Dest.Write(page)
				if err != nil {
					errorExit(fmt.Sprintf(
						"Unexpected error on write: %s", err.Error()))
				}
			}
			if last {
				break
			}
			n++
		}
	}
}

func main() {
	args := args.Get()
	selpg(args)
}
