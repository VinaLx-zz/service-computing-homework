package mux

import (
	"net/http"
)

var mux *http.ServeMux

// Get url multiplexer
func Get() *http.ServeMux {
	if mux == nil {
		mux = makeMux()
	}
	return mux
}

func makeMux() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello world\n"))
	})
	return m
}
