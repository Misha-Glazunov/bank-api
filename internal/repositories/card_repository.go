package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Misha-Glazunov/bank-api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCardNotFound = fmt.Errorf("card not found")
)

type CardRepository interface {
	Create(ctx context.Context, card *models.Card) error
	GetByUserID(ctx context.Context, userID string) ([]*models.Card, error)
}

type PostgresCardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *PostgresCardRepository {
	return &PostgresCardRepository{db: db}
}

func (r *PostgresCardRepository) Create(ctx context.Context, card *models.Card) error {
	// Хеширование CVV
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(card.CVV), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash CVV: %w", err)
	}

	query := `
        INSERT INTO cards (
            user_id,
            number,
            expiry,
            cvv_hash,
            created_at
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		card.UserID,
		card.Number,
		card.Expiry,
		string(cvvHash), // Используем хешированный CVV
		time.Now(),
	).Scan(&card.ID)
}

func (r *PostgresCardRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Card, error) {
	query := `
        SELECT 
            id, 
            user_id, 
            number, 
            expiry, 
            created_at 
        FROM cards 
        WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCardNotFound
		}
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	var cards []*models.Card
	for rows.Next() {
		var card models.Card
		if err := rows.Scan(
			&card.ID,
			&card.UserID,
			&card.Number,
			&card.Expiry,
			&card.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, &card)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return cards, nil
}
