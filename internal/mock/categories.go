package mock

import (
	"bookshop/internal/models"
	"bookshop/internal/repository"
	"errors"
)

func NewCategoryRepository() repository.CategoryRepositoryInterface {
	return &CategoryRepository{}
}

type CategoryRepository struct {
}

func (r *CategoryRepository) Insert(c *models.Category) error {
	switch c.Name {
	case "Duplicate":
		return repository.ErrDuplicateItem
	case "Unexpected":
		return errors.New("empty")
	}
	c.ID = 999
	return nil
}
