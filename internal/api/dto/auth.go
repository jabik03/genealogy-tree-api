package dto

import "time"

// RegisterRequest — запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest — запрос на логин
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse — ответ с токеном
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ProfileResponse — данные профиля пользователя
type ProfileResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateProfileRequest — запрос на обновление профиля
type UpdateProfileRequest struct {
	Email string `json:"email"`
}

// UserResponse — данные пользователя (без пароля!)
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
