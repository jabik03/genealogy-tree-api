package repo

import (
	"GenealogyTree/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateUser создает нового пользователя в БД
func (s *Storage) CreateUser(ctx context.Context, user *models.User) (int, error) {
	query := `
        INSERT INTO users (email, password_hash)
        VALUES ($1, $2)
        RETURNING id, created_at
    `

	err := s.DB.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 - unique_violation (дублирование уникального значения)
			if pgErr.Code == "23505" {
				return 0, ErrEmailAlreadyExists
			}
		}
		return 0, err
	}

	return user.ID, nil
}

// GetUserByEmail находит пользователя по email
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
       SELECT id, email, password_hash, created_at
       FROM users
       WHERE email = $1
    `

	var user models.User
	err := s.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID находит пользователя по ID
func (s *Storage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, created_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := s.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
