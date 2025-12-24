package service

import "GenealogyTree/internal/repo"

type Container struct {
	Person       *PersonService
	Tree         *TreeService
	Relationship *RelationshipService
	Auth         *AuthService
}

func NewContainer(storage *repo.Storage, jwtSecret string) *Container {
	return &Container{
		Person:       NewPersonService(storage),
		Tree:         NewTreeService(storage),
		Relationship: NewRelationshipService(storage),
		Auth:         NewAuthService(storage, jwtSecret),
	}
}
