package repositories

import (
    "context"
    "fmt"
    "database/sql"

    "github.com/Misha-Glazunov/bank-api/internal/models"
)

type TransactionRepository interface {
    Create(ctx context.Context, transaction *models.Transaction) error
    GetByAccountID(ctx context.Context, accountID string) ([]*models.Transaction, error)
}

type PostgresTransactionRepository struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
    return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
    query := `
        INSERT INTO transactions (
            from_account, 
            to_account, 
            amount, 
            type
        ) 
        VALUES ($1, $2, $3, $4)`

    _, err := r.db.ExecContext(ctx, query,
        transaction.FromAccount,
        transaction.ToAccount,
        transaction.Amount,
        transaction.Type,
    )
    return err
}

func (r *PostgresTransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]*models.Transaction, error) {
    query := `
        SELECT id, from_account, to_account, amount, type, created_at 
        FROM transactions 
        WHERE from_account = $1 OR to_account = $1`

    rows, err := r.db.QueryContext(ctx, query, accountID)
    if err != nil {
        return nil, fmt.Errorf("failed to query transactions: %w", err)
    }
    defer rows.Close()

    var transactions []*models.Transaction
    for rows.Next() {
        var t models.Transaction
        err := rows.Scan(
            &t.ID,
            &t.FromAccount,
            &t.ToAccount,
            &t.Amount,
            &t.Type,
            &t.CreatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan transaction: %w", err)
        }
        transactions = append(transactions, &t)
    }
    
    return transactions, nil
}
