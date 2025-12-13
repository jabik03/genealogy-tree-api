package service

import (
	"GenealogyTree/internal/models"
	"GenealogyTree/internal/repo"
	"context"
	"errors"
	"fmt"
)

type RelationshipService struct {
	repo *repo.Storage
}

func NewRelationshipService(storage *repo.Storage) *RelationshipService {
	return &RelationshipService{
		repo: storage,
	}
}

// AddChild связывает существующего ребенка с родителем
func (s *RelationshipService) AddChild(ctx context.Context, parentID, childID int, relType string) (int, error) {
	// Валидация
	if err := s.validateRelationshipType(relType); err != nil {
		return 0, err
	}

	// Получаем родителя и ребенка
	parent, err := s.repo.GetPersonByID(ctx, parentID)
	if err != nil {
		return 0, fmt.Errorf("failed to get parent: %w", err)
	}

	child, err := s.repo.GetPersonByID(ctx, childID)
	if err != nil {
		return 0, fmt.Errorf("failed to get child: %w", err)
	}

	// Проверяем что они в одном дереве
	if parent.TreeID != child.TreeID {
		return 0, errors.New("parent and child must be in the same tree")
	}

	// Проверяем возраст (родитель старше ребенка)
	if err := s.validateAge(parent, child); err != nil {
		return 0, err
	}

	// Проверяем количество родителей у ребенка
	existingParents, err := s.repo.GetParentsByChildID(ctx, childID)
	if err != nil {
		return 0, fmt.Errorf("failed to get existing parents: %w", err)
	}

	if len(existingParents) >= 2 {
		return 0, errors.New("child already has 2 parents")
	}

	// Если уже есть 1 родитель - проверяем пол
	if len(existingParents) == 1 {
		if existingParents[0].IsMale == parent.IsMale {
			return 0, errors.New("parents must be of different gender")
		}
	}

	// Проверяем что такой связи ещё нет
	children, err := s.repo.GetChildrenByParentID(ctx, parentID)
	if err != nil {
		return 0, fmt.Errorf("failed to get existing children: %w", err)
	}

	for _, c := range children {
		if c.ID == childID {
			return 0, errors.New("relationship already exists")
		}
	}

	// Создаём связь
	rel := &models.Relationship{
		ParentID:         parentID,
		ChildID:          childID,
		RelationshipType: relType,
	}

	id, err := s.repo.CreateRelationship(ctx, rel)
	if err != nil {
		return 0, fmt.Errorf("service add child: %w", err)
	}

	return id, nil
}

// AddParent связывает существующего родителя с ребенком
func (s *RelationshipService) AddParent(ctx context.Context, childID, parentID int, relType string) (int, error) {
	// Используем ту же логику что и AddChild (просто меняем местами)
	return s.AddChild(ctx, parentID, childID, relType)
}

// CreateChildAndLink создаёт нового ребенка и связывает с родителем
func (s *RelationshipService) CreateChildAndLink(ctx context.Context, parentID int, child *models.Person, relType string) (int, error) {
	// Валидация
	if err := s.validateRelationshipType(relType); err != nil {
		return 0, err
	}

	// Получаем родителя
	parent, err := s.repo.GetPersonByID(ctx, parentID)
	if err != nil {
		return 0, fmt.Errorf("failed to get parent: %w", err)
	}

	// Устанавливаем TreeID от родителя
	child.TreeID = parent.TreeID

	// Проверяем возраст
	if err := s.validateAge(parent, child); err != nil {
		return 0, err
	}

	// Создаём ребенка
	childID, err := s.repo.CreatePerson(ctx, child)
	if err != nil {
		return 0, fmt.Errorf("failed to create child: %w", err)
	}

	// Создаём связь
	rel := &models.Relationship{
		ParentID:         parentID,
		ChildID:          childID,
		RelationshipType: relType,
	}

	relID, err := s.repo.CreateRelationship(ctx, rel)
	if err != nil {
		return 0, fmt.Errorf("failed to create relationship: %w", err)
	}

	return relID, nil
}

// CreateParentAndLink создаёт нового родителя и связывает с ребенком
func (s *RelationshipService) CreateParentAndLink(ctx context.Context, childID int, parent *models.Person, relType string) (int, error) {
	// Валидация
	if err := s.validateRelationshipType(relType); err != nil {
		return 0, err
	}

	// Получаем ребенка
	child, err := s.repo.GetPersonByID(ctx, childID)
	if err != nil {
		return 0, fmt.Errorf("failed to get child: %w", err)
	}

	// Проверяем количество родителей
	existingParents, err := s.repo.GetParentsByChildID(ctx, childID)
	if err != nil {
		return 0, fmt.Errorf("failed to get existing parents: %w", err)
	}

	if len(existingParents) >= 2 {
		return 0, errors.New("child already has 2 parents")
	}

	// Если есть 1 родитель - проверяем пол
	if len(existingParents) == 1 {
		if existingParents[0].IsMale == parent.IsMale {
			return 0, errors.New("parents must be of different gender")
		}
	}

	// Устанавливаем TreeID от ребенка
	parent.TreeID = child.TreeID

	// Проверяем возраст
	if err := s.validateAge(parent, child); err != nil {
		return 0, err
	}

	// Создаём родителя
	parentID, err := s.repo.CreatePerson(ctx, parent)
	if err != nil {
		return 0, fmt.Errorf("failed to create parent: %w", err)
	}

	// Создаём связь
	rel := &models.Relationship{
		ParentID:         parentID,
		ChildID:          childID,
		RelationshipType: relType,
	}

	relID, err := s.repo.CreateRelationship(ctx, rel)
	if err != nil {
		return 0, fmt.Errorf("failed to create relationship: %w", err)
	}

	return relID, nil
}

// DeleteRelationship удаляет связь
func (s *RelationshipService) DeleteRelationship(ctx context.Context, parentID, childID int) error {
	if parentID <= 0 || childID <= 0 {
		return errors.New("invalid parent or child id")
	}

	if err := s.repo.DeleteRelationship(ctx, parentID, childID); err != nil {
		return fmt.Errorf("service delete relationship: %w", err)
	}

	return nil
}

// GetAvailableChildren получает список доступных детей для персоны
func (s *RelationshipService) GetAvailableChildren(ctx context.Context, parentID int) ([]models.Person, error) {
	if parentID <= 0 {
		return nil, errors.New("invalid parent id")
	}

	// Проверяем что родитель существует
	_, err := s.repo.GetPersonByID(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent: %w", err)
	}

	children, err := s.repo.GetAvailableChildren(ctx, parentID)
	if err != nil {
		return nil, fmt.Errorf("service get available children: %w", err)
	}

	return children, nil
}

// GetAvailableParents получает список доступных родителей для персоны
func (s *RelationshipService) GetAvailableParents(ctx context.Context, childID int) ([]models.Person, error) {
	if childID <= 0 {
		return nil, errors.New("invalid child id")
	}

	// Проверяем что ребенок существует
	_, err := s.repo.GetPersonByID(ctx, childID)
	if err != nil {
		return nil, fmt.Errorf("failed to get child: %w", err)
	}

	parents, err := s.repo.GetAvailableParents(ctx, childID)
	if err != nil {
		return nil, fmt.Errorf("service get available parents: %w", err)
	}

	return parents, nil
}

// validateAge проверяет что родитель старше ребенка
func (s *RelationshipService) validateAge(parent, child *models.Person) error {
	if parent.BirthDate != nil && child.BirthDate != nil {
		if !parent.BirthDate.Before(*child.BirthDate) {
			return errors.New("parent must be born before child")
		}
	}
	return nil
}

// validateRelationshipType проверяет тип связи
func (s *RelationshipService) validateRelationshipType(relType string) error {
	if relType != "biological" && relType != "not_biological" {
		return errors.New("relationship type must be 'biological' or 'not_biological'")
	}
	return nil
}
