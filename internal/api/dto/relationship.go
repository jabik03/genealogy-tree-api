package dto

import "time"

// AddChildRequest — связать существующего ребенка
type AddChildRequest struct {
	ChildID          int    `json:"child_id"`
	RelationshipType string `json:"relationship_type"` // "biological" или "not_biological"
}

// AddParentRequest — связать существующего родителя
type AddParentRequest struct {
	ParentID         int    `json:"parent_id"`
	RelationshipType string `json:"relationship_type"`
}

// CreateChildRequest — создать нового ребенка + связать
type CreateChildRequest struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	BirthDate        *time.Time `json:"birth_date,omitempty"`
	DeathDate        *time.Time `json:"death_date,omitempty"`
	IsMale           bool       `json:"is_male"`
	Biography        string     `json:"biography,omitempty"`
	RelationshipType string     `json:"relationship_type"` // по умолчанию "biological"
}

// CreateParentRequest — создать нового родителя + связать
type CreateParentRequest struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	BirthDate        *time.Time `json:"birth_date,omitempty"`
	DeathDate        *time.Time `json:"death_date,omitempty"`
	IsMale           bool       `json:"is_male"`
	Biography        string     `json:"biography,omitempty"`
	RelationshipType string     `json:"relationship_type"`
}

// RelationshipResponse — ответ после создания связи
type RelationshipResponse struct {
	ID               int    `json:"id"`
	ParentID         int    `json:"parent_id"`
	ChildID          int    `json:"child_id"`
	RelationshipType string `json:"relationship_type"`
}
