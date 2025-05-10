package models

import "time"

type User struct {
    ID           string    `json:"id"`
    Email        string    `json:"email" validate:"required,email"`
    Username     string    `json:"username" validate:"required,alphanum"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}
