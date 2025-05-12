package handlers

import (
	"context"
	"encoding/json"
	"net/http"
		
	"github.com/Misha-Glazunov/bank-api/internal/middleware"
	"github.com/Misha-Glazunov/bank-api/internal/services"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	authService    services.AuthService
	accountService services.AccountService
	cardService    services.CardService
	paymentService services.PaymentService
	cbService      services.CentralBankService
	logger         *logrus.Logger
}

func NewHandlers(
	auth services.AuthService,
	account services.AccountService,
	card services.CardService,
	payment services.PaymentService,
	cb services.CentralBankService,
	logger *logrus.Logger,
) *Handlers {
	return &Handlers{
		authService:    auth,
		accountService: account,
		cardService:    card,
		paymentService: payment,
		cbService:      cb,
		logger:         logger,
	}
}

// Register обработчик регистрации пользователя
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := h.authService.Register(r.Context(), req.Email, req.Username, req.Password); err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, map[string]string{"status": "success"})
}

// Login обработчик аутентификации
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, map[string]string{"token": token})
}

// CreateAccount создание нового счета
func (h *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserIDFromContext(r.Context())
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := h.accountService.CreateAccount(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, account)
}

// CreateCard выпуск новой карты
func (h *Handlers) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	card, err := h.cardService.CreateCard(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, card)
}

// TransferFunds обработчик перевода средств
func (h *Handlers) TransferFunds(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID string  `json:"from_account"`
		ToAccountID   string  `json:"to_account"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if err := h.paymentService.Transfer(r.Context(), req.FromAccountID, req.ToAccountID, req.Amount); err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respondJSON(w, map[string]string{"status": "success"})
}

// Вспомогательные методы

func (h *Handlers) respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *Handlers) respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Errorf("Failed to encode JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handlers) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case services.ErrUserAlreadyExists:
		h.respondError(w, http.StatusConflict, err.Error())
	case services.ErrInvalidCredentials:
		h.respondError(w, http.StatusUnauthorized, err.Error())
	case services.ErrAccountNotFound:
		h.respondError(w, http.StatusNotFound, err.Error())
	case services.ErrInsufficientFunds:
		h.respondError(w, http.StatusBadRequest, err.Error())
	default:
		h.logger.Errorf("Internal server error: %v", err)
		h.respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func getUserIDFromContext(ctx context.Context) (string, error) {
    return middleware.GetUserIDFromContext(ctx)
}

func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
    h.respondJSON(w, map[string]string{"status": "ok"})
}
