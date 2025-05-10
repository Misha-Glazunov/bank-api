package repositories

import (
    "context"
    "database/sql"
    "errors"
    "fmt"

    "github.com/Misha-Glazunov/bank-api/internal/models"
)

var (
    ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    EmailExists(ctx context.Context, email string) (bool, error)
    UsernameExists(ctx context.Context, username string) (bool, error)
}

// Реализация интерфейса с другим именем
type PostgresUserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *PostgresUserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *models.User) error {
    query := `INSERT INTO users (email, username, password_hash)
              VALUES ($1, $2, $3) RETURNING id, created_at`
    return r.db.QueryRowContext(
        ctx,
        query,
        user.Email,
        user.Username,
        user.PasswordHash,
    ).Scan(&user.ID, &user.CreatedAt)
}

// Добавить остальные методы интерфейса
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    query := `SELECT id, email, username, password_hash, created_at 
              FROM users WHERE email = $1`
    row := r.db.QueryRowContext(ctx, query, email)
    
    var user models.User
    err := row.Scan(
        &user.ID,
        &user.Email,
        &user.Username,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("db scan error: %w", err)
    }
    return &user, nil
}

func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    var exists bool
    err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
    return exists, err
}

func (r *PostgresUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
    var exists bool
    err := r.db.QueryRowContext(ctx, query, username).Scan(&exists)
    return exists, err
}
