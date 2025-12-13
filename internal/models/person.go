package models

import "time"

type Person struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	IsMale    bool       `json:"is_male"`
	Biography string     `json:"biography,omitempty"`
	TreeID    int        `json:"tree_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
