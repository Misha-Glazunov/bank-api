package services

import (
    "context"
    "errors"
    
    "github.com/Misha-Glazunov/bank-api/internal/models"
)

// Общие ошибки
var (
    ErrUserAlreadyExists  = errors.New("user already exists")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrAccountNotFound    = errors.New("account not found")
    ErrInsufficientFunds  = errors.New("insufficient funds")
)

type AuthService interface {
    Register(ctx context.Context, email, username, password string) error
    Login(ctx context.Context, email, password string) (string, error)
}

type AccountService interface {
    CreateAccount(ctx context.Context, userID string) (*models.Account, error)
    GetBalance(ctx context.Context, accountID string) (float64, error)
    Deposit(ctx context.Context, accountID string, amount float64) error
    Withdraw(ctx context.Context, accountID string, amount float64) error
}

type CardService interface {
    CreateCard(ctx context.Context, userID string) (*models.Card, error)
}

type CentralBankService interface {
    GetCurrentRate(ctx context.Context) (float64, error)
}

type PaymentService interface {
    Transfer(ctx context.Context, fromAccountID, toAccountID string, amount float64) error
    GetTransactions(ctx context.Context, accountID string) ([]*models.Transaction, error)
}
