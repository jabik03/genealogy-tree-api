package helpers

import (
	"GenealogyTree/internal/api/apierror"
	"net/http"
)

// GetUserIDFromContext извлекает user_id из контекста запроса
func GetUserIDFromContext(r *http.Request) (int, error) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		return 0, apierror.BadRequest("User not authenticated", nil)
	}
	return userID, nil
}
