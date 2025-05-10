package main

import (
    "context"
    "database/sql"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/lib/pq"
    "github.com/sirupsen/logrus"

    "github.com/Misha-Glazunov/bank-api/internal/config"
    "github.com/Misha-Glazunov/bank-api/internal/handlers"
    "github.com/Misha-Glazunov/bank-api/internal/repositories"
    "github.com/Misha-Glazunov/bank-api/internal/routes"
    "github.com/Misha-Glazunov/bank-api/internal/services"
)

func main() {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{})

    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Fatalf("Failed to load config: %v", err)
    }

    connStr := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.DB.Host,
        cfg.DB.Port,
        cfg.DB.User,
        cfg.DB.Password,
        cfg.DB.DBName,
        cfg.DB.SSLMode,
    )

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        logger.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        logger.Fatalf("Database connection failed: %v", err)
    }

    // Инициализация репозиториев
    userRepo := repositories.NewUserRepository(db)
    accountRepo := repositories.NewAccountRepository(db)
    cardRepo := repositories.NewCardRepository(db)
    transactionRepo := repositories.NewTransactionRepository(db)

    // Инициализация сервисов
    authService := services.NewAuthService(userRepo, cfg.JWT.Secret)
    accountService := services.NewAccountService(accountRepo)
    cardService := services.NewCardService(cardRepo)
    paymentService := services.NewPaymentService(accountRepo, transactionRepo)
    centralBankService := services.NewCentralBankService(cfg, logger)

    // Инициализация обработчиков
    h := handlers.NewHandlers(
        authService,
        accountService,
        cardService,
        paymentService,
        centralBankService,
        logger,
    )

    router := routes.NewRouter(h, cfg.JWT.Secret)

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.App.HTTPPort),
        Handler:      router,
        ReadTimeout:  cfg.App.ReadTimeout,
        WriteTimeout: cfg.App.WriteTimeout,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        logger.Infof("Server started on port %d", cfg.App.HTTPPort)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatalf("Server error: %v", err)
        }
    }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    logger.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatalf("Server shutdown failed: %v", err)
    }
    logger.Info("Server exited properly")
}
