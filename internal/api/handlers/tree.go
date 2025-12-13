package handlers

import (
	"GenealogyTree/internal/api/apierror"
	"GenealogyTree/internal/api/dto"
	"GenealogyTree/internal/models"
	"GenealogyTree/internal/repo"
	"GenealogyTree/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TreeHandler struct {
	treeService *service.TreeService
}

func NewTreeHandler(treeService *service.TreeService) *TreeHandler {
	return &TreeHandler{
		treeService: treeService,
	}
}

// CreateTree создаёт новое дерево
func (h *TreeHandler) CreateTree(w http.ResponseWriter, r *http.Request) error {
	var req dto.CreateTreeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// TODO: Когда добавлю JWT, возьму ownerID из токена
	// Пока ownerID = 1 для тестирования
	tree := &models.Tree{
		OwnerID: 1,
		Name:    req.Name,
	}

	id, err := h.treeService.CreateTree(r.Context(), tree)
	if err != nil {
		return apierror.BadRequest("Failed to create tree", err)
	}

	response := dto.TreeResponse{
		ID:        id,
		OwnerID:   tree.OwnerID,
		Name:      tree.Name,
		CreatedAt: tree.CreatedAt,
		UpdatedAt: tree.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// GetTree получает дерево по ID
func (h *TreeHandler) GetTree(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "tree_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	tree, err := h.treeService.GetTreeByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to get tree", err)
	}

	response := dto.TreeResponse{
		ID:        tree.ID,
		OwnerID:   tree.OwnerID,
		Name:      tree.Name,
		CreatedAt: tree.CreatedAt,
		UpdatedAt: tree.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// GetTrees получает все деревья пользователя
func (h *TreeHandler) GetTrees(w http.ResponseWriter, r *http.Request) error {
	// TODO: Когда добавлю JWT ownerID из токена буду брать
	ownerID := 1

	trees, err := h.treeService.GetTreesByOwnerID(r.Context(), ownerID)
	if err != nil {
		return apierror.InternalError("Failed to get trees", err)
	}

	// Формируем список TreeResponse
	treeResponses := make([]dto.TreeResponse, 0, len(trees))
	for _, tree := range trees {
		treeResponses = append(treeResponses, dto.TreeResponse{
			ID:        tree.ID,
			OwnerID:   tree.OwnerID,
			Name:      tree.Name,
			CreatedAt: tree.CreatedAt,
			UpdatedAt: tree.UpdatedAt,
		})
	}
	response := dto.TreeListResponse{
		Trees: treeResponses,
		Total: len(treeResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// UpdateTree обновляет дерево
func (h *TreeHandler) UpdateTree(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "tree_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	var req dto.UpdateTreeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	tree := &models.Tree{
		ID:   id,
		Name: req.Name,
	}

	if err := h.treeService.UpdateTree(r.Context(), tree); err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.BadRequest("Failed to update tree", err)
	}

	// Получаем обновлённое дерево для ответа
	updatedTree, err := h.treeService.GetTreeByID(r.Context(), id)
	if err != nil {
		return apierror.InternalError("Failed to fetch updated tree", err)
	}

	response := dto.TreeResponse{
		ID:        updatedTree.ID,
		OwnerID:   updatedTree.OwnerID,
		Name:      updatedTree.Name,
		CreatedAt: updatedTree.CreatedAt,
		UpdatedAt: updatedTree.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// DeleteTree удаляет дерево
func (h *TreeHandler) DeleteTree(w http.ResponseWriter, r *http.Request) error {
	idStr := chi.URLParam(r, "tree_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	if err := h.treeService.DeleteTree(r.Context(), id); err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to delete tree", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// GetTreeGraph возвращает граф дерева для визуализации
func (h *TreeHandler) GetTreeGraph(w http.ResponseWriter, r *http.Request) error {
	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	persons, relationships, err := h.treeService.GetTreeGraph(r.Context(), treeID)
	if err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to get tree graph", err)
	}

	// Формируем nodes
	nodes := make([]dto.GraphNodeResponse, 0, len(persons))
	for _, person := range persons {
		nodes = append(nodes, dto.GraphNodeResponse{
			ID:        person.ID,
			FirstName: person.FirstName,
			LastName:  person.LastName,
			BirthDate: person.BirthDate,
			DeathDate: person.DeathDate,
			IsMale:    person.IsMale,
		})
	}

	// Формируем edges
	edges := make([]dto.GraphEdgeResponse, 0, len(relationships))
	for _, rel := range relationships {
		edges = append(edges, dto.GraphEdgeResponse{
			ParentID:         rel.ParentID,
			ChildID:          rel.ChildID,
			RelationshipType: rel.RelationshipType,
		})
	}

	response := dto.GraphResponse{
		Nodes: nodes,
		Edges: edges,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
