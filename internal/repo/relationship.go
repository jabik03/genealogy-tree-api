package repo

import (
	"GenealogyTree/internal/models"
	"context"
	"fmt"
)

// CreateRelationship создаёт связь родитель-ребенок
func (s *Storage) CreateRelationship(ctx context.Context, rel *models.Relationship) (int, error) {
	query := `
        INSERT INTO relationships (parent_id, child_id, relationship_type)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err := s.DB.QueryRow(ctx, query,
		rel.ParentID,
		rel.ChildID,
		rel.RelationshipType,
	).Scan(&rel.ID)

	if err != nil {
		return 0, fmt.Errorf("create relationship: %w", err)
	}

	return rel.ID, nil
}

// DeleteRelationship удаляет связь между родителем и ребенком
func (s *Storage) DeleteRelationship(ctx context.Context, parentID, childID int) error {
	query := `DELETE FROM relationships WHERE parent_id = $1 AND child_id = $2`

	commandTag, err := s.DB.Exec(ctx, query, parentID, childID)
	if err != nil {
		return fmt.Errorf("delete relationship: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrRelationshipNotFound
	}

	return nil
}

// GetChildrenByParentID получает всех детей персоны
func (s *Storage) GetChildrenByParentID(ctx context.Context, parentID int) ([]models.Person, error) {
	query := `
        SELECT p.id, p.first_name, p.last_name, p.birth_date, p.death_date,
               p.is_male, p.biography, p.tree_id, p.created_at, p.updated_at
        FROM persons p
        INNER JOIN relationships r ON p.id = r.child_id
        WHERE r.parent_id = $1
    `

	rows, err := s.DB.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("get children by parent: %w", err)
	}
	defer rows.Close()

	var children []models.Person
	for rows.Next() {
		var child models.Person
		if err := rows.Scan(
			&child.ID,
			&child.FirstName,
			&child.LastName,
			&child.BirthDate,
			&child.DeathDate,
			&child.IsMale,
			&child.Biography,
			&child.TreeID,
			&child.CreatedAt,
			&child.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan child: %w", err)
		}
		children = append(children, child)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return children, nil
}

// GetParentsByChildID получает всех родителей персоны
func (s *Storage) GetParentsByChildID(ctx context.Context, childID int) ([]models.Person, error) {
	query := `
        SELECT p.id, p.first_name, p.last_name, p.birth_date, p.death_date,
               p.is_male, p.biography, p.tree_id, p.created_at, p.updated_at
        FROM persons p
        INNER JOIN relationships r ON p.id = r.parent_id
        WHERE r.child_id = $1
    `

	rows, err := s.DB.Query(ctx, query, childID)
	if err != nil {
		return nil, fmt.Errorf("get parents by child: %w", err)
	}
	defer rows.Close()

	var parents []models.Person
	for rows.Next() {
		var parent models.Person
		if err := rows.Scan(
			&parent.ID,
			&parent.FirstName,
			&parent.LastName,
			&parent.BirthDate,
			&parent.DeathDate,
			&parent.IsMale,
			&parent.Biography,
			&parent.TreeID,
			&parent.CreatedAt,
			&parent.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan parent: %w", err)
		}
		parents = append(parents, parent)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return parents, nil
}

// GetAvailableChildren возвращает персон, которые могут быть детьми для данного родителя
func (s *Storage) GetAvailableChildren(ctx context.Context, parentID int) ([]models.Person, error) {
	query := `
        SELECT DISTINCT p.id, p.first_name, p.last_name, p.birth_date, p.is_male, p.tree_id
        FROM persons p
        WHERE p.tree_id = (SELECT tree_id FROM persons WHERE id = $1)
          AND p.id != $1
          AND p.birth_date > (SELECT birth_date FROM persons WHERE id = $1)
          AND p.id NOT IN (
              SELECT child_id FROM relationships WHERE parent_id = $1
          )
          AND (
              SELECT COUNT(*) FROM relationships WHERE child_id = p.id
          ) < 2
        ORDER BY p.birth_date ASC
    `

	rows, err := s.DB.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("get available children: %w", err)
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
			&person.IsMale,
			&person.TreeID,
		); err != nil {
			return nil, fmt.Errorf("scan person: %w", err)
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return persons, nil
}

// GetAvailableParents возвращает персон, которые могут быть родителями для данного ребенка
func (s *Storage) GetAvailableParents(ctx context.Context, childID int) ([]models.Person, error) {
	query := `
        SELECT DISTINCT p.id, p.first_name, p.last_name, p.birth_date, p.is_male, p.tree_id
        FROM persons p
        WHERE p.tree_id = (SELECT tree_id FROM persons WHERE id = $1)
          AND p.id != $1
          AND p.birth_date < (SELECT birth_date FROM persons WHERE id = $1)
          AND p.id NOT IN (
              SELECT parent_id FROM relationships WHERE child_id = $1
          )
          AND (
              (SELECT COUNT(*) FROM relationships WHERE child_id = $1) = 0
              OR p.is_male != (
                  SELECT is_male 
                  FROM persons 
                  WHERE id = (SELECT parent_id FROM relationships WHERE child_id = $1 LIMIT 1)
              )
          )
        ORDER BY p.birth_date DESC
    `

	rows, err := s.DB.Query(ctx, query, childID)
	if err != nil {
		return nil, fmt.Errorf("get available parents: %w", err)
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
			&person.IsMale,
			&person.TreeID,
		); err != nil {
			return nil, fmt.Errorf("scan person: %w", err)
		}
		persons = append(persons, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return persons, nil
}
