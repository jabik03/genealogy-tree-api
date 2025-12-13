package service

import (
	"GenealogyTree/internal/models"
	"GenealogyTree/internal/repo"
	"context"
	"errors"
	"fmt"
	"time"
)

// PersonService содержит бизнес-логику для работы с персонами
type PersonService struct {
	repo *repo.Storage // Доступ к репозиторию
}

// NewPersonService создаёт новый сервис
func NewPersonService(storage *repo.Storage) *PersonService {
	return &PersonService{
		repo: storage,
	}
}

// CreatePerson создаёт новую персону с валидацией
func (s *PersonService) CreatePerson(ctx context.Context, p *models.Person) (int, error) {
	// 1. Валидация структуры
	if err := s.validatePerson(p); err != nil {
		return 0, err
	}

	if p.BirthDate != nil && p.BirthDate.After(time.Now()) {
		return 0, errors.New("birth date cannot be in the future")
	}

	if p.BirthDate != nil && p.DeathDate != nil {
		if p.DeathDate.Before(*p.BirthDate) {
			return 0, errors.New("death date cannot be before birth date")
		}
	}

	// 4. Вызываем репозиторий для создания записи
	id, err := s.repo.CreatePerson(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("service create person: %w", err)
	}

	return id, nil
}

// GetPersonByID получает персону по ID
func (s *PersonService) GetPersonByID(ctx context.Context, id int) (*models.Person, error) {
	if id <= 0 {
		return nil, errors.New("invalid person id")
	}

	person, err := s.repo.GetPersonByID(ctx, id)
	if err != nil {
		// Пробрасываем ошибку репозитория (включая ErrPersonNotFound)
		return nil, fmt.Errorf("service get person: %w", err)
	}

	return person, nil
}

// GetPersonsByTreeID получает все персоны в конкретном дереве
func (s *PersonService) GetPersonsByTreeID(ctx context.Context, treeID int) ([]models.Person, error) {
	if treeID <= 0 {
		return nil, errors.New("invalid tree id")
	}

	persons, err := s.repo.GetPersonsByTreeID(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("service get persons by tree: %w", err)
	}

	return persons, nil
}

// UpdatePerson обновляет персону с валидацией
func (s *PersonService) UpdatePerson(ctx context.Context, p *models.Person) error {
	// Валидация
	if err := s.validatePerson(p); err != nil {
		return err
	}

	if p.ID <= 0 {
		return errors.New("invalid person id")
	}

	if p.BirthDate != nil && p.BirthDate.After(time.Now()) {
		return errors.New("birth date cannot be in the future")
	}

	if p.BirthDate != nil && p.DeathDate != nil {
		if p.DeathDate.Before(*p.BirthDate) {
			return errors.New("death date cannot be before birth date")
		}
	}

	// Вызываем репозиторий
	if err := s.repo.UpdatePerson(ctx, p); err != nil {
		return fmt.Errorf("service update person: %w", err)
	}

	return nil
}

// DeletePerson удаляет персону
func (s *PersonService) DeletePerson(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid person id")
	}

	if err := s.repo.DeletePerson(ctx, id); err != nil {
		return fmt.Errorf("service delete person: %w", err)
	}

	return nil
}

func (s *PersonService) validatePerson(p *models.Person) error {
	if p.FirstName == "" {
		return errors.New("first name is required")
	}
	if p.LastName == "" {
		return errors.New("last name is required")
	}
	if p.TreeID <= 0 {
		return errors.New("tree id is required")
	}

	return nil
}
