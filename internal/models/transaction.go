package models

import "time"

type Transaction struct {
    ID          string    `json:"id" db:"id"`
    FromAccount string    `json:"from_account" db:"from_account"`
    ToAccount   string    `json:"to_account" db:"to_account"`
    Amount      float64   `json:"amount" db:"amount"`
    Type        string    `json:"type" db:"type"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
