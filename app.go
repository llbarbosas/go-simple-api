package main

import (
	"net/http"
)

func main() {
	a := &app{
		UserHandler: new(UserHandler),
	}
	http.ListenAndServe(":8000", a)
}

type app struct {
	UserHandler *UserHandler
}

func (a *app) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	res.Header().Set("Content-Type", "application/json")

	switch head {
	case "user":
		a.UserHandler.ServeHTTP(res, req)
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}
