package server

import (
	"args"
	"fmt"
	"io"
	"log"
	"net/http"
	"server/mux"
)

func getLoggerFromWriter(writer io.Writer) *log.Logger {
	return log.New(writer, "http", log.LstdFlags)
}

// Serve at specific host and port
func Serve(args *args.Args) {
	server := http.Server{
		Addr:     fmt.Sprintf("%s:%d", args.Host, args.Port),
		Handler:  mux.Get(),
		ErrorLog: getLoggerFromWriter(args.Log),
	}
	fmt.Fprintf(args.Log, "hello world serving at %s:%d\n", args.Host, args.Port)
	server.ListenAndServe()
}
