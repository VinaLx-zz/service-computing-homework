package server

import (
	"args"
	"net/http"
	"router"

	"github.com/urfave/negroni"
)

func initStatic(path string) *negroni.Static {
	static := negroni.NewStatic(http.Dir(path))
	static.Prefix = "/static"
	return static
}

func initRecovery() *negroni.Recovery {
	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	recovery.Formatter = new(negroni.TextPanicFormatter)
	return recovery
}

// New creates a new server
func New(a *args.Args) *negroni.Negroni {
	server := negroni.New(negroni.NewLogger())
	server.Use(initStatic(a.StaticPath))
	server.Use(initRecovery())
	server.UseHandler(router.Get())
	return server
}
