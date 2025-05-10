package services

import (
    "context"
    "fmt"

    "github.com/Misha-Glazunov/bank-api/internal/models"
    "github.com/Misha-Glazunov/bank-api/internal/repositories"
)

type paymentServiceImpl struct {
    accountRepo     repositories.AccountRepository
    transactionRepo repositories.TransactionRepository
}

func NewPaymentService(
    accountRepo repositories.AccountRepository,
    transactionRepo repositories.TransactionRepository,
) PaymentService {
    return &paymentServiceImpl{
        accountRepo:     accountRepo,
        transactionRepo: transactionRepo,
    }
}

func (s *paymentServiceImpl) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount float64) error {
    if err := s.accountRepo.UpdateBalance(ctx, fromAccountID, -amount); err != nil {
        return fmt.Errorf("withdrawal failed: %w", err)
    }

    if err := s.accountRepo.UpdateBalance(ctx, toAccountID, amount); err != nil {
        s.accountRepo.UpdateBalance(ctx, fromAccountID, amount) // Rollback
        return fmt.Errorf("deposit failed: %w", err)
    }
    return nil
}

func (s *paymentServiceImpl) GetTransactions(ctx context.Context, accountID string) ([]*models.Transaction, error) {
    return s.transactionRepo.GetByAccountID(ctx, accountID)
}
