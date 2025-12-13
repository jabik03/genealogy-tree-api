package service

import "GenealogyTree/internal/repo"

type Container struct {
	Person       *PersonService
	Tree         *TreeService
	Relationship *RelationshipService
	// User   *UserService      // Добавим позже
}

func NewContainer(storage *repo.Storage) *Container {
	return &Container{
		Person:       NewPersonService(storage),
		Tree:         NewTreeService(storage),
		Relationship: NewRelationshipService(storage),
		// User:   NewUserService(storage),
	}
}
