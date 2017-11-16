package args

import "github.com/ogier/pflag"

// Args specifies the command line args of the server
type Args struct {
	Host       string
	Port       uint16
	StaticPath string
}

var arg *Args

// Get args
func Get() *Args {
	if arg == nil {
		arg = parseArgs()
	}
	return arg
}

func parseArgs() *Args {
	host := pflag.StringP("host", "h", "localhost", "host of the server serves")
	port := pflag.Uint16P("port", "p", 8080, "port of the server serves")
	static := pflag.StringP("static", "s", "static", "the root of static files")
	pflag.Parse()
	return &Args{
		Host:       *host,
		Port:       *port,
		StaticPath: *static,
	}
}
