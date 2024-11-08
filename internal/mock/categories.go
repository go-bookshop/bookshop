package mock

import (
	"bookshop/internal/data"
	"errors"
)

func NewCategoriesRepository() data.CategoriesRepositoryInterface {
	return &CategoriesRepository{}
}

type CategoriesRepository struct {
}

func (r *CategoriesRepository) Insert(c *data.Category) error {
	switch c.Name {
	case "Duplicate":
		return data.ErrDuplicateItem
	case "Unexpected":
		return errors.New("empty")
	}
	c.ID = 999
	return nil
}
