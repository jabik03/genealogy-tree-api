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

type RelationshipHandler struct {
	relationshipService *service.RelationshipService
}

func NewRelationshipHandler(relationshipService *service.RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{
		relationshipService: relationshipService,
	}
}

// AddChild связывает существующего ребенка с родителем
func (h *RelationshipHandler) AddChild(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	var req dto.AddChildRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// По умолчанию biological
	if req.RelationshipType == "" {
		req.RelationshipType = "biological"
	}

	relID, err := h.relationshipService.AddChild(r.Context(), personID, req.ChildID, req.RelationshipType)
	if err != nil {
		return apierror.BadRequest("Failed to add child", err)
	}

	response := dto.RelationshipResponse{
		ID:               relID,
		ParentID:         personID,
		ChildID:          req.ChildID,
		RelationshipType: req.RelationshipType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// AddParent связывает существующего родителя с ребенком
func (h *RelationshipHandler) AddParent(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	var req dto.AddParentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// По умолчанию biological
	if req.RelationshipType == "" {
		req.RelationshipType = "biological"
	}

	relID, err := h.relationshipService.AddParent(r.Context(), personID, req.ParentID, req.RelationshipType)
	if err != nil {
		return apierror.BadRequest("Failed to add parent", err)
	}

	response := dto.RelationshipResponse{
		ID:               relID,
		ParentID:         req.ParentID,
		ChildID:          personID,
		RelationshipType: req.RelationshipType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// CreateChildAndLink создаёт нового ребенка и связывает с родителем
func (h *RelationshipHandler) CreateChildAndLink(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	var req dto.CreateChildRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// По умолчанию biological
	if req.RelationshipType == "" {
		req.RelationshipType = "biological"
	}

	child := &models.Person{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		DeathDate: req.DeathDate,
		IsMale:    req.IsMale,
		Biography: req.Biography,
	}

	relID, err := h.relationshipService.CreateChildAndLink(r.Context(), personID, child, req.RelationshipType)
	if err != nil {
		return apierror.BadRequest("Failed to create child and link", err)
	}

	response := dto.RelationshipResponse{
		ID:               relID,
		ParentID:         personID,
		ChildID:          child.ID,
		RelationshipType: req.RelationshipType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// CreateParentAndLink создаёт нового родителя и связывает с ребенком
func (h *RelationshipHandler) CreateParentAndLink(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	var req dto.CreateParentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	// По умолчанию biological
	if req.RelationshipType == "" {
		req.RelationshipType = "biological"
	}

	parent := &models.Person{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		DeathDate: req.DeathDate,
		IsMale:    req.IsMale,
		Biography: req.Biography,
	}

	relID, err := h.relationshipService.CreateParentAndLink(r.Context(), personID, parent, req.RelationshipType)
	if err != nil {
		return apierror.BadRequest("Failed to create parent and link", err)
	}

	response := dto.RelationshipResponse{
		ID:               relID,
		ParentID:         parent.ID,
		ChildID:          personID,
		RelationshipType: req.RelationshipType,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// RemoveChild удаляет связь родитель-ребенок
func (h *RelationshipHandler) RemoveChild(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	childIDStr := chi.URLParam(r, "child_id")
	childID, err := strconv.Atoi(childIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid child ID format", err)
	}

	if err := h.relationshipService.DeleteRelationship(r.Context(), personID, childID); err != nil {
		if errors.Is(err, repo.ErrRelationshipNotFound) {
			return apierror.NotFound("Relationship not found", err)
		}
		return apierror.InternalError("Failed to remove child", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// RemoveParent удаляет связь ребенок-родитель
func (h *RelationshipHandler) RemoveParent(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	parentIDStr := chi.URLParam(r, "parent_id")
	parentID, err := strconv.Atoi(parentIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid parent ID format", err)
	}

	// Удаляем связь (parent -> child)
	if err := h.relationshipService.DeleteRelationship(r.Context(), parentID, personID); err != nil {
		if errors.Is(err, repo.ErrRelationshipNotFound) {
			return apierror.NotFound("Relationship not found", err)
		}
		return apierror.InternalError("Failed to remove parent", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// GetAvailableChildren возвращает список доступных детей для персоны
func (h *RelationshipHandler) GetAvailableChildren(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	children, err := h.relationshipService.GetAvailableChildren(r.Context(), personID)
	if err != nil {
		if errors.Is(err, repo.ErrPersonNotFound) {
			return apierror.NotFound("Person not found", err)
		}
		return apierror.InternalError("Failed to get available children", err)
	}

	// Формируем краткий ответ
	childrenResponses := make([]dto.PersonBriefResponse, 0, len(children))
	for _, child := range children {
		childrenResponses = append(childrenResponses, dto.PersonBriefResponse{
			ID:        child.ID,
			FirstName: child.FirstName,
			LastName:  child.LastName,
			BirthDate: child.BirthDate,
			IsMale:    child.IsMale,
		})
	}

	response := dto.PersonBriefListResponse{
		Persons: childrenResponses,
		Total:   len(childrenResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// GetAvailableParents возвращает список доступных родителей для персоны
func (h *RelationshipHandler) GetAvailableParents(w http.ResponseWriter, r *http.Request) error {
	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	parents, err := h.relationshipService.GetAvailableParents(r.Context(), personID)
	if err != nil {
		if errors.Is(err, repo.ErrPersonNotFound) {
			return apierror.NotFound("Person not found", err)
		}
		return apierror.InternalError("Failed to get available parents", err)
	}

	// Формируем краткий ответ
	parentsResponses := make([]dto.PersonBriefResponse, 0, len(parents))
	for _, parent := range parents {
		parentsResponses = append(parentsResponses, dto.PersonBriefResponse{
			ID:        parent.ID,
			FirstName: parent.FirstName,
			LastName:  parent.LastName,
			BirthDate: parent.BirthDate,
			IsMale:    parent.IsMale,
		})
	}

	response := dto.PersonBriefListResponse{
		Persons: parentsResponses,
		Total:   len(parentsResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}
