package middleware

import (
	"GenealogyTree/internal/api/apierror"
	"GenealogyTree/internal/api/dto"
	"GenealogyTree/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// AuthMiddleware проверяет JWT токен и добавляет user_id в контекст
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, apierror.BadRequest("Authorization header required", nil))
				return
			}

			// 2. Проверяем формат: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, apierror.BadRequest("Invalid authorization header format", nil))
				return
			}

			tokenString := parts[1]

			// 3. Валидируем токен и получаем user_id
			userID, err := authService.ValidateToken(tokenString)
			if err != nil {
				respondWithError(w, apierror.BadRequest("Invalid or expired token", err))
				return
			}

			// 4. Добавляем user_id в контекст запроса
			ctx := context.WithValue(r.Context(), "user_id", userID)

			// 5. Передаём управление следующему handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// respondWithError отправляет JSON-ответ с ошибкой
func respondWithError(w http.ResponseWriter, apiErr *apierror.APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.StatusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: apiErr.Message})
}
