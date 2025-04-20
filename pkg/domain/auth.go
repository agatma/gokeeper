package domain

import "github.com/google/uuid"

type InUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           uuid.UUID `json:"uid"`
	Login        string    `json:"login"`
	PasswordHash []byte    `json:"-"`
}
