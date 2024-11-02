package mock

import (
	"bookshop/internal/data"
	"errors"
)

func NewAuthorRepository() data.AuthorRepositoryInterface {
	return &AuthorRepository{}
}

type AuthorRepository struct {
}

func (r *AuthorRepository) Insert(a *data.Author) error {
	switch a.Name {
	case "Duplicate":
		return data.ErrDuplicateAuthorName
	case "Unexpected":
		return errors.New("empty")
	}
	a.ID = 999
	return nil
}
