package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	users []*User
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

type Role int8

const (
	UserRole  Role = 0
	AdminRole Role = 1
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password []byte    `json:"password"`
	Role     `json:"role"`
}

type PlainUser struct {
	Name     string
	Email    string
	Password string
}

func NewUser(data PlainUser) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 2)

	if err != nil {
		return nil, err
	}

	matched, err := regexp.MatchString("@golang.com$", data.Email)

	if err != nil {
		return nil, err
	}

	var role Role

	if matched {
		role = AdminRole
	} else {
		role = UserRole
	}

	return &User{uuid.New(), data.Name, data.Email, passwordHash, role}, nil
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

	h.users = append(h.users, user)

	js, err := json.Marshal(user)

	res.Write(js)
}
