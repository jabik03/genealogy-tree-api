package app

import (
	"GenealogyTree/internal/api"
	"GenealogyTree/internal/config"
	_ "GenealogyTree/internal/logger"
	"GenealogyTree/internal/repo"
	"GenealogyTree/internal/service"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Info("⚠️ No .env file found")
	}
}

func Run(ctx context.Context) error {
	conf := config.NewConfig()
	port := fmt.Sprintf(":%s", conf.ApiConf.Port)

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	storage, err := repo.NewDB(dbCtx, conf)
	if err != nil {
		slog.Error("❌ Failed to connect to database", "error", err)
		return err
	}
	defer storage.Close()

	slog.Info("✅ Connected to database", "db", conf.Database.Name)

	// Создаём контейнер со ВСЕМИ сервисами
	services := service.NewContainer(storage, conf.JWT.SecretKey)
	slog.Info("✅ Services initialized")

	// Передаём контейнер в роутер
	router := api.NewRouter(services)
	server := api.NewServer(port, router.Mux)

	if err := server.Start(ctx); err != nil {
		slog.Error("🔴 Server stopped with error", "error", err)
		return err
	}

	slog.Info("🟢 Server exited gracefully")
	return nil
}
