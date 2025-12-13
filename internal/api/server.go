package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		slog.Info("🚀 Server starting (http://localhost:8080)", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("💥 Server error", "err", err)
		}
	}()

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Info("🛑 Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("⚠️ Server shutdown error", "err", err)
		return fmt.Errorf("server shutdown error: %w", err)
	}

	slog.Info("✅ Server stopped gracefully")
	return nil
}
