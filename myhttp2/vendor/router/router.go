package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"table"

	"github.com/gorilla/mux"
)

var r *mux.Router

// Get router
func Get() *mux.Router {
	if r == nil {
		r = initRouter()
	}
	return r
}

type jsRespond struct {
	Status  bool
	Message string
}

func testjs(w http.ResponseWriter, r *http.Request) {
	respond := jsRespond{Status: true}
	if r.Method == "GET" {
		respond.Message = "get ok"
	} else if r.Method == "POST" {
		respond.Message = "post ok"
	} else {
		respond.Status = false
		respond.Message = "invalid method"
	}
	json.NewEncoder(w).Encode(respond)
}

func notfound(w http.ResponseWriter, r *http.Request) {
	panic(fmt.Sprintf(
		"Error: unknown path: %s, maybe it's under development?", r.URL.Path))
}

func formtable(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	table.Render(r.Form, w)
}

func initRouter() *mux.Router {
	r = mux.NewRouter()
	r.HandleFunc("/testjs", testjs)
	r.HandleFunc("/formtable", formtable)
	r.PathPrefix("/").HandlerFunc(notfound)
	return r
}
