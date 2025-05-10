// internal/repositories/account_repository.go
package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Misha-Glazunov/bank-api/internal/models"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Account, error)
	UpdateBalance(ctx context.Context, accountID string, amount float64) error
}

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Create(ctx context.Context, account *models.Account) error {
	query := `
		INSERT INTO accounts (
			user_id, 
			balance, 
			currency, 
			created_at
		) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		account.UserID,
		account.Balance,
		account.Currency,
		time.Now(),
	).Scan(&account.ID, &account.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, id string) (*models.Account, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			balance, 
			currency, 
			created_at 
		FROM accounts 
		WHERE id = $1`

	var account models.Account
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &account, nil
}

func (r *PostgresAccountRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Account, error) {
	query := `
		SELECT 
			id, 
			user_id, 
			balance, 
			currency, 
			created_at 
		FROM accounts 
		WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Balance,
			&account.Currency,
			&account.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, &account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return accounts, nil
}

func (r *PostgresAccountRepository) UpdateBalance(ctx context.Context, accountID string, amount float64) error {
	query := `
		UPDATE accounts 
		SET balance = balance + $1 
		WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, amount, accountID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAccountNotFound
	}

	return nil
}
