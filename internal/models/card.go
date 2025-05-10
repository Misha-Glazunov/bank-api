package models

import "time"

type Card struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Number    string    `json:"number" db:"number"`
    Expiry    string    `json:"expiry" db:"expiry"`
    CVV       string    `json:"-"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
