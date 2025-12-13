package api

import (
	"GenealogyTree/internal/api/apierror"
	"GenealogyTree/internal/api/dto"
	"GenealogyTree/internal/api/handlers"
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
}

func NewRouter(services *service.Container) *Router {
	r := &Router{
		Mux:                 chi.NewRouter(),
		services:            services,
		personHandler:       handlers.NewPersonHandler(services.Person),
		treeHandler:         handlers.NewTreeHandler(services.Tree),
		relationshipHandler: handlers.NewRelationshipHandler(services.Relationship),
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
	r.Mux.Get("/ping", r.handler(handlers.PingHandler))

	// CRUD для trees
	r.Mux.Get("/api/trees", r.handler(r.treeHandler.GetTrees))
	r.Mux.Post("/api/trees", r.handler(r.treeHandler.CreateTree))
	r.Mux.Get("/api/trees/{tree_id}", r.handler(r.treeHandler.GetTree))
	r.Mux.Put("/api/trees/{tree_id}", r.handler(r.treeHandler.UpdateTree))
	r.Mux.Delete("/api/trees/{tree_id}", r.handler(r.treeHandler.DeleteTree))

	// CRUD для persons (вложил в trees)
	r.Mux.Get("/api/trees/{tree_id}/persons", r.handler(r.personHandler.GetPersons))
	r.Mux.Post("/api/trees/{tree_id}/persons", r.handler(r.personHandler.CreatePerson))
	r.Mux.Get("/api/trees/{tree_id}/persons/{person_id}", r.handler(r.personHandler.GetPerson))
	r.Mux.Put("/api/trees/{tree_id}/persons/{person_id}", r.handler(r.personHandler.UpdatePerson))
	r.Mux.Delete("/api/trees/{tree_id}/persons/{person_id}", r.handler(r.personHandler.DeletePerson))

	// Работа с Relationships
	// Relationships - связать существующих
	r.Mux.Post("/api/persons/{person_id}/children", r.handler(r.relationshipHandler.AddChild))
	r.Mux.Post("/api/persons/{person_id}/parents", r.handler(r.relationshipHandler.AddParent))

	// Relationships - создать нового + связать
	r.Mux.Post("/api/persons/{person_id}/children/new", r.handler(r.relationshipHandler.CreateChildAndLink))
	r.Mux.Post("/api/persons/{person_id}/parents/new", r.handler(r.relationshipHandler.CreateParentAndLink))

	// Relationships - удалить связь
	r.Mux.Delete("/api/persons/{person_id}/children/{child_id}", r.handler(r.relationshipHandler.RemoveChild))
	r.Mux.Delete("/api/persons/{person_id}/parents/{parent_id}", r.handler(r.relationshipHandler.RemoveParent))

	// Available children/parents
	r.Mux.Get("/api/persons/{person_id}/available-children", r.handler(r.relationshipHandler.GetAvailableChildren))
	r.Mux.Get("/api/persons/{person_id}/available-parents", r.handler(r.relationshipHandler.GetAvailableParents))

	// Graph endpoint
	r.Mux.Get("/api/trees/{tree_id}/graph", r.handler(r.treeHandler.GetTreeGraph))
}

func (r *Router) handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			r.handleError(w, req, err)
		}
	}
}

func (r *Router) handleError(w http.ResponseWriter, req *http.Request, err error) {
	// Пытаюсь извлечь APIError
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		// Это наша кастомная ошибка
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

	// Неизвестная ошибка - возвращаем 500
	slog.Error("unexpected error",
		"method", req.Method,
		"path", req.URL.Path,
		"err", err,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: "Internal Server Error"})
}
