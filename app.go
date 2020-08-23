package main

import (
	"database/sql"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		log.Fatal(err)
	}

	a := &App{
		UserHandler: NewUserHandler(db),
	}

	err = http.ListenAndServe(":8000", a)

	if err != nil {
		log.Fatal(err)
	}
}

type App struct {
	UserHandler *UserHandler
}

func (a *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
