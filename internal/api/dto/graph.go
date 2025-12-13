package dto

import "time"

// GraphNodeResponse — узел графа (персона)
type GraphNodeResponse struct {
	ID        int        `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	DeathDate *time.Time `json:"death_date,omitempty"`
	IsMale    bool       `json:"is_male"`
}

// GraphEdgeResponse — связь (ребро графа)
type GraphEdgeResponse struct {
	ParentID         int    `json:"parent_id"`
	ChildID          int    `json:"child_id"`
	RelationshipType string `json:"relationship_type"`
}

// GraphResponse — полный граф дерева
type GraphResponse struct {
	Nodes []GraphNodeResponse `json:"nodes"`
	Edges []GraphEdgeResponse `json:"edges"`
}
