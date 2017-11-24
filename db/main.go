package main

import (
	"server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	server.Start()
}
