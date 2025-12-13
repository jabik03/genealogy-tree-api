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

type PersonHandler struct {
	personService *service.PersonService
}

func NewPersonHandler(personService *service.PersonService) *PersonHandler {
	return &PersonHandler{
		personService: personService,
	}
}

// GetPersons получает все персоны в дереве
func (h *PersonHandler) GetPersons(w http.ResponseWriter, r *http.Request) error {
	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	persons, err := h.personService.GetPersonsByTreeID(r.Context(), treeID)
	if err != nil {
		return apierror.InternalError("Failed to get persons", err)
	}

	// Формируем список PersonResponse
	personResponses := make([]dto.PersonResponse, 0, len(persons))
	for _, person := range persons {
		personResponses = append(personResponses, dto.PersonResponse{
			ID:        person.ID,
			FirstName: person.FirstName,
			LastName:  person.LastName,
			BirthDate: person.BirthDate,
			DeathDate: person.DeathDate,
			IsMale:    person.IsMale,
			Biography: person.Biography,
			TreeID:    person.TreeID,
			CreatedAt: person.CreatedAt,
			UpdatedAt: person.UpdatedAt,
		})
	}

	response := dto.PersonListResponse{
		Persons: personResponses,
		Total:   len(personResponses),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// CreatePerson создаёт персону в дереве
func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) error {
	treeIDStr := chi.URLParam(r, "tree_id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	var req dto.CreatePersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	person := &models.Person{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		DeathDate: req.DeathDate,
		IsMale:    req.IsMale,
		Biography: req.Biography,
		TreeID:    treeID, // Берём из URL
	}

	id, err := h.personService.CreatePerson(r.Context(), person)
	if err != nil {
		return apierror.BadRequest("Failed to create person", err)
	}

	response := dto.PersonResponse{
		ID:        id,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		BirthDate: person.BirthDate,
		DeathDate: person.DeathDate,
		IsMale:    person.IsMale,
		Biography: person.Biography,
		TreeID:    person.TreeID,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

// GetPerson получает персону по ID
func (h *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) error {
	treeIDStr := chi.URLParam(r, "tree_id")
	_, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	person, err := h.personService.GetPersonByID(r.Context(), personID)
	if err != nil {
		if errors.Is(err, repo.ErrPersonNotFound) {
			return apierror.NotFound("Person not found", err)
		}
		return apierror.InternalError("Failed to get person", err)
	}

	response := dto.PersonResponse{
		ID:        person.ID,
		FirstName: person.FirstName,
		LastName:  person.LastName,
		BirthDate: person.BirthDate,
		DeathDate: person.DeathDate,
		IsMale:    person.IsMale,
		Biography: person.Biography,
		TreeID:    person.TreeID,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// UpdatePerson обновляет персону
func (h *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) error {
	treeIDStr := chi.URLParam(r, "tree_id")
	_, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	var req dto.UpdatePersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return apierror.BadRequest("Invalid JSON", err)
	}

	person := &models.Person{
		ID:        personID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: req.BirthDate,
		DeathDate: req.DeathDate,
		IsMale:    req.IsMale,
		Biography: req.Biography,
		// TreeID не трогаем - он не меняется в UPDATE
	}

	if err := h.personService.UpdatePerson(r.Context(), person); err != nil {
		if errors.Is(err, repo.ErrPersonNotFound) {
			return apierror.NotFound("Person not found", err)
		}
		return apierror.BadRequest("Failed to update person", err)
	}

	updatedPerson, err := h.personService.GetPersonByID(r.Context(), personID)
	if err != nil {
		return apierror.InternalError("Failed to fetch updated person", err)
	}

	response := dto.PersonResponse{
		ID:        updatedPerson.ID,
		FirstName: updatedPerson.FirstName,
		LastName:  updatedPerson.LastName,
		BirthDate: updatedPerson.BirthDate,
		DeathDate: updatedPerson.DeathDate,
		IsMale:    updatedPerson.IsMale,
		Biography: updatedPerson.Biography,
		TreeID:    updatedPerson.TreeID,
		CreatedAt: updatedPerson.CreatedAt,
		UpdatedAt: updatedPerson.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

// DeletePerson удаляет персону
func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) error {
	treeIDStr := chi.URLParam(r, "tree_id")
	_, err := strconv.Atoi(treeIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid tree ID format", err)
	}

	personIDStr := chi.URLParam(r, "person_id")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		return apierror.BadRequest("Invalid person ID format", err)
	}

	if err := h.personService.DeletePerson(r.Context(), personID); err != nil {
		if errors.Is(err, repo.ErrPersonNotFound) {
			return apierror.NotFound("Person not found", err)
		}
		return apierror.InternalError("Failed to delete person", err)
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
