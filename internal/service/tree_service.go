package service

import (
	"GenealogyTree/internal/models"
	"GenealogyTree/internal/repo"
	"context"
	"errors"
	"fmt"
)

type TreeService struct {
	repo *repo.Storage
}

func NewTreeService(storage *repo.Storage) *TreeService {
	return &TreeService{
		repo: storage,
	}
}

// CreateTree создаёт новое дерево с валидацией
func (s *TreeService) CreateTree(ctx context.Context, t *models.Tree) (int, error) {
	if err := s.validateTree(t); err != nil {
		return 0, err
	}

	if t.OwnerID <= 0 {
		return 0, errors.New("owner id is required")
	}

	id, err := s.repo.CreateTree(ctx, t)
	if err != nil {
		return 0, fmt.Errorf("service create tree: %w", err)
	}

	return id, nil
}

// GetTreeByID получает дерево по ID
func (s *TreeService) GetTreeByID(ctx context.Context, id int) (*models.Tree, error) {
	if id <= 0 {
		return nil, errors.New("invalid tree id")
	}

	tree, err := s.repo.GetTreeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service get tree: %w", err)
	}

	return tree, nil
}

func (s *TreeService) GetTreesByOwnerID(ctx context.Context, ownerID int) ([]models.Tree, error) {
	if ownerID <= 0 {
		return nil, errors.New("invalid owner id")
	}

	trees, err := s.repo.GetTreesByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("service get trees by owner: %w", err)
	}

	return trees, nil
}

func (s *TreeService) UpdateTree(ctx context.Context, t *models.Tree) error {
	if err := s.validateTree(t); err != nil {
		return err
	}

	if t.ID <= 0 {
		return errors.New("invalid tree id")
	}

	if err := s.repo.UpdateTree(ctx, t); err != nil {
		return fmt.Errorf("service update tree: %w", err)
	}

	return nil
}

func (s *TreeService) DeleteTree(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid tree id")
	}

	if err := s.repo.DeleteTree(ctx, id); err != nil {
		return fmt.Errorf("service delete tree: %w", err)
	}

	return nil
}

// GetTreeGraph получает граф дерева (персоны + связи)
func (s *TreeService) GetTreeGraph(ctx context.Context, treeID int) ([]models.Person, []models.Relationship, error) {
	if treeID <= 0 {
		return nil, nil, errors.New("invalid tree id")
	}

	_, err := s.repo.GetTreeByID(ctx, treeID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tree: %w", err)
	}

	persons, relationships, err := s.repo.GetTreeGraph(ctx, treeID)
	if err != nil {
		return nil, nil, fmt.Errorf("service get tree graph: %w", err)
	}

	return persons, relationships, nil
}

func (s *TreeService) validateTree(t *models.Tree) error {
	if t.Name == "" {
		return errors.New("tree name is required")
	}

	if len(t.Name) > 255 {
		return errors.New("tree name is too long (max 255 characters)")
	}

	return nil
}
