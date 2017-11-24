package main

import (
	"args"
	"io/ioutil"
	"log"
	"server"

	_ "github.com/mattn/go-sqlite3"
)

func setLog() {
	if !*args.Log {
		log.SetOutput(ioutil.Discard)
	}
}

func main() {
	setLog()
	server.Start()
}
