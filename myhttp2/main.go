package main

import (
	"args"
	"fmt"
	"server"
)

func main() {
	a := args.Get()
	sv := server.New(a)
	sv.Run(fmt.Sprintf("%s:%d", a.Host, a.Port))
}
