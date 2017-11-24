package router

import (
	"args"
	"database/ormdao"
	"database/sqldao"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"user"

	"github.com/gorilla/mux"
)

// Response ..
type response struct {
	OK   bool
	Data interface{}
}

func jsonResponse(w http.ResponseWriter, resp *response) {
	encoder := json.NewEncoder(w)
	encoder.Encode(resp)
}

func success(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	jsonResponse(w, &response{
		OK:   true,
		Data: data,
	})
}
func failure(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	log.Printf("INFO: %s", message)
	jsonResponse(w, &response{
		OK:   false,
		Data: message,
	})
}

func illegalForm(w http.ResponseWriter, form url.Values, params ...string) bool {
	for _, p := range params {
		if form[p] == nil {
			failure(
				w, http.StatusBadRequest,
				fmt.Sprintf("illegal input: require field: %s", p))
			return true
		}
	}
	return false
}

func getDao() (user.Dao, error) {
	if *args.ORM {
		return ormdao.Get()
	}
	return sqldao.Get()
}

func adduser(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if illegalForm(w, req.Form, "username", "password") {
		return
	}
	u := user.NewUser(req.Form["username"][0], req.Form["password"][0])
	dao, err := getDao()
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf(
		"INFO: adding user, username '%s', password '%s'",
		u.Username, u.Password)
	err = dao.StoreUser(u)
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
	} else {
		success(w, "add user success")
	}
}
func getuser(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if illegalForm(w, req.Form, "userid") {
		return
	}
	uid, err := strconv.ParseUint(req.Form["userid"][0], 10, 64)
	if err != nil {
		failure(w, http.StatusBadRequest, err.Error())
		return
	}
	dao, err := getDao()
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	u, err := dao.GetUser(uid)
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
	} else {
		success(w, u)
	}
}
func getallusers(w http.ResponseWriter, req *http.Request) {
	dao, err := getDao()
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
		return
	}
	us, err := dao.GetAllUsers()
	if err != nil {
		failure(w, http.StatusInternalServerError, err.Error())
	} else {
		success(w, us)
	}
}

// New router
func New() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/adduser", adduser).Methods("POST", "GET")
	r.HandleFunc("/getuser", getuser).Methods("POST", "GET")
	r.HandleFunc("/getallusers", getallusers).Methods("GET")
	r.HandleFunc("/", http.NotFound)
	return r
}
