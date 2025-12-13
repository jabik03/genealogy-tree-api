package repo

import (
	"GenealogyTree/internal/config"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg *config.Config) (*Storage, error) {

	dsn := cfg.Database.BuildDSN()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("new pgxpool: %v", err)
	}

	// Проверим соединение
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgxpool: %v", err)
	}

	slog.Debug("connected to postgres")

	return &Storage{DB: pool}, nil
}

func (s *Storage) Close() {
	if s.DB != nil {
		s.DB.Close()
		slog.Debug("database connection pool closed") // Правильное сообщение
	} else {
		slog.Debug("database pool was already nil")
	}
}
