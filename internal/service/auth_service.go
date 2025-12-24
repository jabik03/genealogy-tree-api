package service

import (
	"GenealogyTree/internal/models"
	"GenealogyTree/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repo.Storage
	secretKey string
}

func NewAuthService(storage *repo.Storage, secretKey string) *AuthService {
	return &AuthService{
		repo:      storage,
		secretKey: secretKey,
	}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(ctx context.Context, email, password string) (*models.User, error) {
	// Валидация
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаём пользователя
	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repo.ErrEmailAlreadyExists) {
			return nil, errors.New("email already registered")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login проверяет credentials и возвращает JWT токен
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// Валидация
	if email == "" {
		return "", errors.New("email is required")
	}
	if password == "" {
		return "", errors.New("password is required")
	}

	// Ищем пользователя по email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Сравниваем пароль с хэшем
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Создаём JWT токен
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

// generateToken создаёт JWT токен для пользователя
func (s *AuthService) generateToken(userID int) (string, error) {
	// Создаём claims (полезная нагрузка)
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Токен живёт 24 часа
		"iat":     time.Now().Unix(),                     // Время создания
	}

	// Создаём токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет JWT токен и возвращает user_id
func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Проверяем что токен валидный
	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Извлекаем user_id
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return int(userID), nil
}

func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID получает пользователя по ID (для профиля)
func (s *AuthService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
