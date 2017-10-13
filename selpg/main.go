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

func printOne(src *args.ReadSrc, a *args.Args) {
	if len(src.Name) != 0 {
		a.Dest.Write([]byte(fmt.Sprintf("%s:\n", src.Name)))
	}
	n := 0
	for page := range a.Pager(src.Reader) {
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

func selpg(a *args.Args) {
	for source := range a.Sources {
		printOne(source, a)
		if source.Next != nil {
			source.Next <- true
		}
	}
}

func main() {
	args := args.Get()
	selpg(args)
}
