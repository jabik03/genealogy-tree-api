package dto

import "time"

// CreateTreeRequest — данные для создания дерева
type CreateTreeRequest struct {
	Name string `json:"name"`
	// OwnerID будем брать из JWT токена (после авторизации)
	// Пока можно временно передавать или захардкодить
}

// UpdateTreeRequest — данные для обновления дерева
type UpdateTreeRequest struct {
	Name string `json:"name"`
	// ID берётся из URL параметра
	// OwnerID не меняется
}

// TreeResponse — данные дерева в ответе
type TreeResponse struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TreeListResponse — для списка деревьев
type TreeListResponse struct {
	Trees []TreeResponse `json:"trees"`
	Total int            `json:"total"`
}
