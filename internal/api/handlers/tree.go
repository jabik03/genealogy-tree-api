package handlers

import (
	"GenealogyTree/internal/api/apierror"
	"GenealogyTree/internal/api/dto"
	"GenealogyTree/internal/api/helpers"
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

// GetTrees получает все деревья текущего пользователя
func (h *TreeHandler) GetTrees(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	trees, err := h.treeService.GetTreesByOwnerID(r.Context(), userID)
	if err != nil {
		return apierror.InternalError("Failed to get trees", err)
	}

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

// CreateTree создаёт новое дерево для текущего пользователя
func (h *TreeHandler) CreateTree(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	var req dto.CreateTreeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	tree := &models.Tree{
		OwnerID: userID,
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

// GetTree получает дерево по ID (с проверкой владельца)
func (h *TreeHandler) GetTree(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	tree, err := h.treeService.GetTreeByID(r.Context(), treeID)
	if err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to get tree", err)
	}

	// Проверяем что пользователь владелец этого дерева
	if tree.OwnerID != userID {
		return apierror.NotFound("Tree not found", nil)
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

// UpdateTree обновляет дерево (с проверкой владельца)
func (h *TreeHandler) UpdateTree(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	// Проверяем что дерево существует и принадлежит пользователю
	existingTree, err := h.treeService.GetTreeByID(r.Context(), treeID)
	if err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to get tree", err)
	}

	if existingTree.OwnerID != userID {
		return apierror.NotFound("Tree not found", nil)
	}

	var req dto.UpdateTreeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	tree := &models.Tree{
		ID:   treeID,
		Name: req.Name,
	}

	if err := h.treeService.UpdateTree(r.Context(), tree); err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.BadRequest("Failed to update tree", err)
	}

	response := dto.TreeResponse{
		ID:        tree.ID,
		OwnerID:   existingTree.OwnerID,
		Name:      tree.Name,
		CreatedAt: existingTree.CreatedAt,
		UpdatedAt: tree.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// DeleteTree удаляет дерево (с проверкой владельца)
func (h *TreeHandler) DeleteTree(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	// Проверяем что дерево существует и принадлежит пользователю
	tree, err := h.treeService.GetTreeByID(r.Context(), treeID)
	if err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to get tree", err)
	}

	if tree.OwnerID != userID {
		return apierror.NotFound("Tree not found", nil)
	}

	if err := h.treeService.DeleteTree(r.Context(), treeID); err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to delete tree", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{
		"message": "Tree deleted successfully",
	})
}

// GetTreeGraph возвращает граф дерева для визуализации (с проверкой владельца)
func (h *TreeHandler) GetTreeGraph(w http.ResponseWriter, r *http.Request) error {
	userID, err := helpers.GetUserIDFromContext(r)
	if err != nil {
		return err
	}

	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	// Проверяем что дерево существует и принадлежит пользователю
	tree, err := h.treeService.GetTreeByID(r.Context(), treeID)
	if err != nil {
		if errors.Is(err, repo.ErrTreeNotFound) {
			return apierror.NotFound("Tree not found", err)
		}
		return apierror.InternalError("Failed to get tree", err)
	}

	if tree.OwnerID != userID {
		return apierror.NotFound("Tree not found", nil)
	}

	persons, relationships, err := h.treeService.GetTreeGraph(r.Context(), treeID)
	if err != nil {
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
