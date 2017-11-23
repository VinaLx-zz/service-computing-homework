package server

import (
	"args"
	"fmt"
	"router"

	"github.com/urfave/negroni"
)

// Start the server
func Start() {
	s := new()
	s.UseHandler(router.New())
	s.Run(fmt.Sprintf("%s:%d", *args.Host, *args.Port))
}

func new() *negroni.Negroni {
	return negroni.Classic()
}
