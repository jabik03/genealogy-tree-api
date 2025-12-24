package handlers

import (
	"GenealogyTree/internal/api/apierror"
	"GenealogyTree/internal/api/dto"
	"GenealogyTree/internal/api/helpers"
	"GenealogyTree/internal/repo"
	"GenealogyTree/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register регистрирует нового пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// Вызываем сервис
	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, repo.ErrEmailAlreadyExists) || err.Error() == "email already registered" {
			return apierror.BadRequest("Email already registered", err)
		}
		return apierror.BadRequest(err.Error(), err)
	}

	// Генерируем токен для нового пользователя
	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		return apierror.InternalError("Failed to generate token", err)
	}

	// Формируем ответ
	response := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// Login авторизует пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// Вызываем сервис
	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return apierror.BadRequest("Invalid email or password", err)
		}
		return apierror.InternalError("Failed to login", err)
	}

	// Получаем данные пользователя для ответа
	user, err := h.authService.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		return apierror.InternalError("Failed to get user", err)
	}

	// Формируем ответ
	response := dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// Logout — клиентский logout (токен удаляется на фронтенде)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	userID, ok := r.Context().Value("user_id").(int)
	if ok {
		slog.Info("User logged out", "user_id", userID)
	}

	response := map[string]string{
		"message": "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// GetProfile возвращает профиль текущего пользователя
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	user, err := h.authService.GetUserByID(r.Context(), userID)
	if err != nil {
		return apierror.InternalError("Failed to get profile", err)
	}

	response := dto.ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
