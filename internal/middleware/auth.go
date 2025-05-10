// internal/middleware/auth.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Ключ для сохранения userID в контексте
type contextKey string

const userIDKey contextKey = "userID"

// AuthMiddleware создает middleware для JWT-аутентификации
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Извлечение токена из заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Проверка формата заголовка
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			tokenString := tokenParts[1]

			// Парсинг токена
			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(jwtSecret), nil
				},
			)

			// Обработка ошибок парсинга
			if err != nil || !token.Valid {
				respondError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Извлечение userID из claims
			userID := claims.Subject
			if userID == "" {
				respondError(w, http.StatusUnauthorized, "Malformed token")
				return
			}

			// Добавление userID в контекст
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext извлекает userID из контекста
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// Вспомогательная функция для отправки ошибок
func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
