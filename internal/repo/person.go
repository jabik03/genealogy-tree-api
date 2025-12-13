package repo

import (
	"GenealogyTree/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreatePerson(ctx context.Context, p *models.Person) (int, error) {
	query := `
		INSERT INTO persons (first_name, last_name, birth_date, death_date, is_male, biography, tree_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
	`

	err := s.DB.QueryRow(ctx, query,
		p.FirstName,
		p.LastName,
		p.BirthDate,
		p.DeathDate,
		p.IsMale,
		p.Biography,
		p.TreeID,
	).Scan(&p.ID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return 0, fmt.Errorf("create person: %w", err)
	}

	return p.ID, nil
}

func (s *Storage) GetPersonByID(ctx context.Context, id int) (*models.Person, error) {
	query := `
		SELECT id, first_name, last_name, birth_date, death_date, 
		       is_male, biography, tree_id, created_at, updated_at
        FROM persons
        WHERE id = $1
	`
	var person models.Person

	row := s.DB.QueryRow(ctx, query, id)

	err := row.Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.BirthDate,
		&person.DeathDate,
		&person.IsMale,
		&person.Biography,
		&person.TreeID,
		&person.CreatedAt,
		&person.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPersonNotFound
		}
		return nil, fmt.Errorf("get person by id: %w", err)
	}

	return &person, nil
}

func (s *Storage) GetPersonsByTreeID(ctx context.Context, treeID int) ([]models.Person, error) {
	query := `
        SELECT id, first_name, last_name, birth_date, death_date, 
               is_male, biography, tree_id, created_at, updated_at
        FROM persons
        WHERE tree_id = $1
        ORDER BY created_at DESC
    `

	rows, err := s.DB.Query(ctx, query, treeID)
	if err != nil {
		return nil, fmt.Errorf("get persons by tree: %w", err)
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
			&person.Biography,
			&person.TreeID,
			&person.CreatedAt,
			&person.UpdatedAt,
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

func (s *Storage) UpdatePerson(ctx context.Context, p *models.Person) error {
	query := `
       UPDATE persons 
       SET first_name = $1, 
           last_name = $2, 
           birth_date = $3, 
           death_date = $4, 
           is_male = $5, 
           biography = $6,
           updated_at = NOW()
       WHERE id = $7
    `

	commandTag, err := s.DB.Exec(ctx, query,
		p.FirstName, // $1
		p.LastName,  // $2
		p.BirthDate, // $3
		p.DeathDate, // $4
		p.IsMale,    // $5
		p.Biography, // $6
		p.ID,        // $7
	)

	if err != nil {
		return fmt.Errorf("update person: %w", err)
	}

	// была ли обновлена хотя бы одна строка
	if commandTag.RowsAffected() == 0 {
		return ErrPersonNotFound
	}

	return nil
}

func (s *Storage) DeletePerson(ctx context.Context, id int) error {
	query := `DELETE FROM persons WHERE id = $1`

	commandTag, err := s.DB.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete person: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrPersonNotFound
	}
	return nil
}
