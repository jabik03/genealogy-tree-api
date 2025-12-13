package repo

import (
	"GenealogyTree/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateTree(ctx context.Context, t *models.Tree) (int, error) {
	query := `
		INSERT INTO trees (owner_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := s.DB.QueryRow(ctx, query,
		t.OwnerID,
		t.Name,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return 0, fmt.Errorf("create tree: %w", err)
	}

	return t.ID, nil
}

func (s *Storage) GetTreeByID(ctx context.Context, id int) (*models.Tree, error) {
	query := `
		SELECT id, owner_id, name, created_at, updated_at
		FROM trees
		WHERE id = $1
	`
	var tree models.Tree

	err := s.DB.QueryRow(ctx, query, id).Scan(
		&tree.ID,
		&tree.OwnerID,
		&tree.Name,
		&tree.CreatedAt,
		&tree.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTreeNotFound
		}
		return nil, fmt.Errorf("get tree by id: %w", err)
	}

	return &tree, nil
}

func (s *Storage) GetTreesByOwnerID(ctx context.Context, ownerID int) ([]models.Tree, error) {
	query := `
		SELECT id, owner_id, name, created_at, updated_at
        FROM trees
        WHERE owner_id = $1
        ORDER BY created_at DESC
	`
	rows, err := s.DB.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("get trees by owner: %w", err)
	}
	defer rows.Close()

	var trees []models.Tree
	for rows.Next() {
		var tree models.Tree
		if err := rows.Scan(
			&tree.ID,
			&tree.OwnerID,
			&tree.Name,
			&tree.CreatedAt,
			&tree.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan tree: %w", err)
		}
		trees = append(trees, tree)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return trees, nil
}

func (s *Storage) UpdateTree(ctx context.Context, t *models.Tree) error {
	query := `
		UPDATE trees
        SET name = $1,
            updated_at = NOW()
		WHERE id = $2
	`

	commandTag, err := s.DB.Exec(ctx, query,
		t.Name,
		t.ID,
	)

	if err != nil {
		return fmt.Errorf("update tree: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrTreeNotFound
	}

	return nil
}

func (s *Storage) DeleteTree(ctx context.Context, id int) error {
	query := `DELETE FROM trees WHERE id = $1`

	commandTag, err := s.DB.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete tree: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrTreeNotFound
	}

	return nil
}

// GetTreeGraph получает все персоны и связи для визуализации дерева
func (s *Storage) GetTreeGraph(ctx context.Context, treeID int) ([]models.Person, []models.Relationship, error) {
	// Получаем ТОЛЬКО нужные поля для графа
	personsQuery := `
        SELECT id, first_name, last_name, birth_date, death_date, is_male
        FROM persons
        WHERE tree_id = $1
        ORDER BY birth_date ASC
    `

	rows, err := s.DB.Query(ctx, personsQuery, treeID)
	if err != nil {
		return nil, nil, fmt.Errorf("get persons for graph: %w", err)
	}
	defer rows.Close()

	var persons []models.Person
	for rows.Next() {
		var person models.Person
		if err := rows.Scan(
			&person.ID,
			&person.FirstName,
			&person.LastName,
			&person.BirthDate,
			&person.DeathDate,
			&person.IsMale,
		); err != nil {
			return nil, nil, fmt.Errorf("scan person: %w", err)
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("rows error: %w", err)
	}

	// Получаем все связи для персон в этом дереве
	relationshipsQuery := `
        SELECT r.id, r.parent_id, r.child_id, r.relationship_type
        FROM relationships r
        INNER JOIN persons p ON r.parent_id = p.id
        WHERE p.tree_id = $1
    `

	relRows, err := s.DB.Query(ctx, relationshipsQuery, treeID)
	if err != nil {
		return nil, nil, fmt.Errorf("get relationships for graph: %w", err)
	}
	defer relRows.Close()

	var relationships []models.Relationship
	for relRows.Next() {
		var rel models.Relationship
		if err := relRows.Scan(
			&rel.ID,
			&rel.ParentID,
			&rel.ChildID,
			&rel.RelationshipType,
		); err != nil {
			return nil, nil, fmt.Errorf("scan relationship: %w", err)
		}
		relationships = append(relationships, rel)
	}

	if err := relRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("relationships rows error: %w", err)
	}

	return persons, relationships, nil
}
