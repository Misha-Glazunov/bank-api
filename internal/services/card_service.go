package services

import (
    "context"
    
    "github.com/Misha-Glazunov/bank-api/internal/models"
    "github.com/Misha-Glazunov/bank-api/internal/repositories"
)

type cardServiceImpl struct {
    repo repositories.CardRepository
}

func NewCardService(repo repositories.CardRepository) CardService {
    return &cardServiceImpl{repo: repo}
}

func (s *cardServiceImpl) CreateCard(ctx context.Context, userID string) (*models.Card, error) {
    card := &models.Card{
        UserID: userID,
        Number: "4111111111111111",
        Expiry: "12/30",
    }
    err := s.repo.Create(ctx, card)
    return card, err
}
