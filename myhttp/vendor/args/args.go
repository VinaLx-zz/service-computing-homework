package args

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Args pass to the server
type Args struct {
	Host string
	Port uint
	Log  io.WriteCloser
}

var args *Args

// Get command line arguments
func Get() *Args {
	if args == nil {
		args = parseArgs()
	}
	return args
}

func getLogWriter(file string) io.WriteCloser {
	if file == "stdout" || file == "" {
		return os.Stdout
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0744)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	return f
}

func parseArgs() *Args {
	port := flag.Uint("port", 8080, "port of the server listening to")
	host := flag.String("host", "localhost", "host of the server listening to")
	logfile := flag.String("log", "stdout", "file output of log")
	flag.Parse()
	return &Args{
		Host: *host,
		Port: *port,
		Log:  getLogWriter(*logfile),
	}
}
