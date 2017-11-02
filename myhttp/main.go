package main

import (
	"args"
	"server"
)

func main() {
	a := args.Get()
	server.Serve(a)
}
