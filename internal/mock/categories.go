package mock

import (
	"bookshop/internal/data"
	"errors"
)

func NewCategoryRepository() data.CategoryRepositoryInterface {
	return &CategoryRepository{}
}

type CategoryRepository struct {
}

func (r *CategoryRepository) Insert(c *data.Category) error {
	switch c.Name {
	case "Duplicate":
		return data.ErrDuplicateItem
	case "Unexpected":
		return errors.New("empty")
	}
	c.ID = 999
	return nil
}
