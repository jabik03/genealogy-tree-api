package repo

import "errors"

// Sentinel errors для репозитория
var (
	ErrTreeNotFound         = errors.New("tree not found")
	ErrPersonNotFound       = errors.New("person not found")
	ErrRelationshipNotFound = errors.New("relationship not found")
)
