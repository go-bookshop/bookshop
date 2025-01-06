package models

import (
	"bookshop/internal/validator"
	"fmt"
	"strings"
	"time"
)

const (
	CategoryNameMaxLength        int = 100
	CategoryDescriptionMaxLength int = 1000
)

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func ValidateCategory(v *validator.Validator, c *Category) {
	v.Check(strings.TrimSpace(c.Name) != "", "name", "must be provided")
	v.Check(len(c.Name) <= CategoryNameMaxLength, "name", fmt.Sprintf("must be less than %d bytes long", CategoryNameMaxLength))

	v.Check(strings.TrimSpace(c.Description) != "", "description", "must be provided")
	v.Check(len(c.Description) <= CategoryDescriptionMaxLength, "description", fmt.Sprintf("must be less than %d bytes long", CategoryDescriptionMaxLength))
}
