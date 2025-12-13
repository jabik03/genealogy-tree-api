package dto

import "time"

// CreatePersonRequest — данные для создания персоны
type CreatePersonRequest struct {
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	IsMale    bool       `json:"is_male"`
	Biography string     `json:"biography,omitempty"`
}

// UpdatePersonRequest — данные для обновления персоны
type UpdatePersonRequest struct {
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	IsMale    bool       `json:"is_male"`
	Biography string     `json:"biography,omitempty"`
}

// PersonResponse — данные персоны в ответе
type PersonResponse struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	IsMale    bool       `json:"is_male"`
	Biography string     `json:"biography,omitempty"`
	TreeID    int        `json:"tree_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// PersonListResponse — для списка персон
type PersonListResponse struct {
	Persons []PersonResponse `json:"persons"`
	Total   int              `json:"total"`
}

// PersonBriefResponse — краткая информация о персоне для списков
type PersonBriefResponse struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	IsMale    bool       `json:"is_male"`
}

// PersonBriefListResponse — список кратких данных
type PersonBriefListResponse struct {
	Persons []PersonBriefResponse `json:"persons"`
	Total   int                   `json:"total"`
}

// ErrorResponse — структура для ошибок
type ErrorResponse struct {
	Error string `json:"error"`
}
