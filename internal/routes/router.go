package routes

import (
    "github.com/gorilla/mux"
    "github.com/Misha-Glazunov/bank-api/internal/handlers"
    "github.com/Misha-Glazunov/bank-api/internal/middleware"
)

func NewRouter(h *handlers.Handlers, jwtSecret string) *mux.Router {
    r := mux.NewRouter()
    
    r.HandleFunc("/healthcheck", h.HealthCheck).Methods("GET")
    r.HandleFunc("/register", h.Register).Methods("POST")
    r.HandleFunc("/login", h.Login).Methods("POST")
    
    authRouter := r.PathPrefix("/").Subrouter()
    authRouter.Use(middleware.AuthMiddleware(jwtSecret))
    
    authRouter.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
    authRouter.HandleFunc("/transfer", h.TransferFunds).Methods("POST")
    
    return r
}
