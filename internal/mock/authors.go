package mock

import (
	"bookshop/internal/models"
	"bookshop/internal/repository"
	"errors"
)

func NewAuthorRepository() repository.AuthorRepositoryInterface {
	return &AuthorRepository{}
}

type AuthorRepository struct {
}

func (r *AuthorRepository) Insert(a *models.Author) error {
	switch a.Name {
	case "Unexpected":
		return errors.New("empty")
	}
	a.ID = 999
	return nil
}
