package services

import (
    "context"
    "fmt"

    "github.com/Misha-Glazunov/bank-api/internal/models"
    "github.com/Misha-Glazunov/bank-api/internal/repositories"
)

type accountServiceImpl struct {
    repo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
    return &accountServiceImpl{repo: repo}
}

func (s *accountServiceImpl) CreateAccount(ctx context.Context, userID string) (*models.Account, error) {
    account := &models.Account{
        UserID:  userID,
        Balance: 0.0,
    }
    err := s.repo.Create(ctx, account)
    return account, err
}

func (s *accountServiceImpl) GetBalance(ctx context.Context, accountID string) (float64, error) {
    account, err := s.repo.GetByID(ctx, accountID)
    if err != nil {
        return 0, ErrAccountNotFound
    }
    return account.Balance, nil
}

func (s *accountServiceImpl) Deposit(ctx context.Context, accountID string, amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("invalid deposit amount")
    }
    return s.repo.UpdateBalance(ctx, accountID, amount)
}

func (s *accountServiceImpl) Withdraw(ctx context.Context, accountID string, amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("invalid withdrawal amount")
    }
    
    balance, err := s.GetBalance(ctx, accountID)
    if err != nil {
        return err
    }
    
    if balance < amount {
        return ErrInsufficientFunds
    }
    
    return s.repo.UpdateBalance(ctx, accountID, -amount)
}
