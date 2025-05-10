package services

import (
    "context"
    "fmt"
    "time"

    "github.com/Misha-Glazunov/bank-api/internal/models"
    "github.com/Misha-Glazunov/bank-api/internal/repositories"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

// Удален ненужный импорт "errors"

type authServiceImpl struct {
    userRepo  repositories.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
    return &authServiceImpl{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
    }
}

func (s *authServiceImpl) Register(ctx context.Context, email, username, password string) error {
    exists, err := s.userRepo.EmailExists(ctx, email)
    if err != nil {
        return fmt.Errorf("email check failed: %w", err)
    }
    if exists {
        return ErrUserAlreadyExists // Используем ошибку из interfaces.go
    }

    exists, err = s.userRepo.UsernameExists(ctx, username)
    if err != nil {
        return fmt.Errorf("username check failed: %w", err)
    }
    if exists {
        return ErrUserAlreadyExists
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("password hashing failed: %w", err)
    }

    user := &models.User{
        Email:        email,
        Username:     username,
        PasswordHash: string(hashedPassword),
    }

    return s.userRepo.Create(ctx, user)
}

func (s *authServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return "", ErrInvalidCredentials // Используем ошибку из interfaces.go
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", ErrInvalidCredentials
    }

    return GenerateJWTToken(user.ID, s.jwtSecret)
}

func GenerateJWTToken(userID string, secret string) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   userID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
