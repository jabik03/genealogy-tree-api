package models

type Relationship struct {
	ID               int    `json:"id"`
	ParentID         int    `json:"parent_id"`
	ChildID          int    `json:"child_id"`
	RelationshipType string `json:"relationship_type"`
}
