package models

import "time"

type Tree struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
