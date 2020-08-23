package main

import (
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
