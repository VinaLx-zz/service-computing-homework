package args

import (
	"github.com/ogier/pflag"
)

// Host ..
var Host *string

// Port ..
var Port *uint16

// ORM ..
var ORM *bool

// DB path
var DB *string

// Log enable
var Log *bool

func init() {
	Host = pflag.StringP("host", "h", "localhost", "host of the server")
	Port = pflag.Uint16P("port", "p", 8080, "port of the server")
	ORM = pflag.BoolP("orm", "o", false, "use ORM for database")
	DB = pflag.StringP("db", "d", "./users.db", "database used by sqlite")
	Log = pflag.BoolP("log", "l", false, "enable logging mode")

	pflag.Parse()
}
