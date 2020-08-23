package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type UserHandler struct {
	*sql.DB
}

func (h *UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	match := Matcher(req.Method, head)

	switch {
	case match("GET", ""):
		h.Create(res, req)
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}

func NewUserHandler(db *sql.DB) *UserHandler {
	handler := new(UserHandler)

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS USERS (id text, name text, email text, password text, role integer)")

	if err != nil {
		log.Fatal(err)
	}

	handler.DB = db
	return handler
}

// Create a new user
func (h *UserHandler) Create(res http.ResponseWriter, req *http.Request) {
	var u PlainUser

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&u)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := NewUser(u)

	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	js, err := json.Marshal(user)

	res.Write(js)
}
