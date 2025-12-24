package api

import (
	"GenealogyTree/internal/api/apierror"
	"GenealogyTree/internal/api/dto"
	"GenealogyTree/internal/api/handlers"
	appMiddleware "GenealogyTree/internal/api/middleware" // ← Импорт нашего middleware
	"GenealogyTree/internal/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

type Router struct {
	Mux                 *chi.Mux
	services            *service.Container
	personHandler       *handlers.PersonHandler
	treeHandler         *handlers.TreeHandler
	relationshipHandler *handlers.RelationshipHandler
	authHandler         *handlers.AuthHandler
}

func NewRouter(services *service.Container) *Router {
	r := &Router{
		Mux:                 chi.NewRouter(),
		services:            services,
		personHandler:       handlers.NewPersonHandler(services.Person),
		treeHandler:         handlers.NewTreeHandler(services.Tree),
		relationshipHandler: handlers.NewRelationshipHandler(services.Relationship),
		authHandler:         handlers.NewAuthHandler(services.Auth),
	}

	r.initMiddleware()
	r.initRoutes()

	return r
}

func (r *Router) initMiddleware() {
	r.Mux.Use(middleware.Logger)
	r.Mux.Use(middleware.Recoverer)
}

func (r *Router) initRoutes() {
	// ========================================
	// Тестовый endpoints
	// ========================================
	r.Mux.Get("/ping", r.handler(handlers.PingHandler))

	// ========================================
	// Публичные endpoints (БЕЗ токена)
	// ========================================
	r.Mux.Post("/api/auth/register", r.handler(r.authHandler.Register))
	r.Mux.Post("/api/auth/login", r.handler(r.authHandler.Login))

	// ========================================
	// Защищённые endpoints (ТРЕБУЮТ токен)
	// ========================================
	r.Mux.Group(func(protected chi.Router) {
		// Применяем Auth Middleware ко всем роутам в этой группе
		protected.Use(appMiddleware.AuthMiddleware(r.services.Auth))

		// Auth
		protected.Post("/api/auth/logout", r.handler(r.authHandler.Logout))
		protected.Get("/api/profile", r.handler(r.authHandler.GetProfile))

		// Trees
		protected.Get("/api/trees", r.handler(r.treeHandler.GetTrees))
		protected.Post("/api/trees", r.handler(r.treeHandler.CreateTree))
		protected.Get("/api/trees/{tree_id}", r.handler(r.treeHandler.GetTree))
		protected.Put("/api/trees/{tree_id}", r.handler(r.treeHandler.UpdateTree))
		protected.Delete("/api/trees/{tree_id}", r.handler(r.treeHandler.DeleteTree))

		// Persons
		protected.Get("/api/trees/{tree_id}/persons", r.handler(r.personHandler.GetPersons))
		protected.Post("/api/trees/{tree_id}/persons", r.handler(r.personHandler.CreatePerson))
		protected.Get("/api/trees/{tree_id}/persons/{person_id}", r.handler(r.personHandler.GetPerson))
		protected.Put("/api/trees/{tree_id}/persons/{person_id}", r.handler(r.personHandler.UpdatePerson))
		protected.Delete("/api/trees/{tree_id}/persons/{person_id}", r.handler(r.personHandler.DeletePerson))

		// Relationships
		protected.Post("/api/persons/{person_id}/children", r.handler(r.relationshipHandler.AddChild))
		protected.Post("/api/persons/{person_id}/parents", r.handler(r.relationshipHandler.AddParent))
		protected.Post("/api/persons/{person_id}/children/new", r.handler(r.relationshipHandler.CreateChildAndLink))
		protected.Post("/api/persons/{person_id}/parents/new", r.handler(r.relationshipHandler.CreateParentAndLink))
		protected.Delete("/api/persons/{person_id}/children/{child_id}", r.handler(r.relationshipHandler.RemoveChild))
		protected.Delete("/api/persons/{person_id}/parents/{parent_id}", r.handler(r.relationshipHandler.RemoveParent))

		// Available
		protected.Get("/api/persons/{person_id}/available-children", r.handler(r.relationshipHandler.GetAvailableChildren))
		protected.Get("/api/persons/{person_id}/available-parents", r.handler(r.relationshipHandler.GetAvailableParents))

		// Graph
		protected.Get("/api/trees/{tree_id}/graph", r.handler(r.treeHandler.GetTreeGraph))
	})
}

func (r *Router) handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			r.handleError(w, req, err)
		}
	}
}

func (r *Router) handleError(w http.ResponseWriter, req *http.Request, err error) {
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		slog.Error("API error",
			"method", req.Method,
			"path", req.URL.Path,
			"status", apiErr.StatusCode,
			"message", apiErr.Message,
			"err", apiErr.Err,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.StatusCode)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Error: apiErr.Message})
		return
	}

	slog.Error("unexpected error",
		"method", req.Method,
		"path", req.URL.Path,
		"err", err,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Internal Server Error"})
}
