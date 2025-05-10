package services

import (
    "context"
    "net/http"

    "github.com/Misha-Glazunov/bank-api/internal/config"
    "github.com/sirupsen/logrus"
)

type centralBankServiceImpl struct {
    client *http.Client
    config *config.CentralCBConfig
    logger *logrus.Logger
}

func NewCentralBankService(cfg *config.Config, logger *logrus.Logger) CentralBankService {
    return &centralBankServiceImpl{
        client: &http.Client{Timeout: cfg.CentralCB.Timeout},
        config: &cfg.CentralCB,
        logger: logger,
    }
}

func (s *centralBankServiceImpl) GetCurrentRate(ctx context.Context) (float64, error) {
    // Временная заглушка
    return 15.0, nil
}
