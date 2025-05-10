package repositories

import (
    "context"
    "database/sql"
    "time"

    "github.com/Misha-Glazunov/bank-api/internal/models"
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
        "", // Реализовать хеширование
        time.Now(),
    ).Scan(&card.ID)
}

func (r *PostgresCardRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Card, error) {
    // Реализация получения карт пользователя
    return nil, nil
}
