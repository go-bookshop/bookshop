package models

import (
	"bookshop/internal/validator"
	"fmt"
	"strings"
	"time"
)

const (
	AuthorNameMaxLength = 1000
	AuthorBioMaxLength  = 2000
)

type Author struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func ValidateAuthor(v *validator.Validator, a *Author) {
	v.Check(strings.TrimSpace(a.Name) != "", "name", "must be provided")
	v.Check(len(a.Name) <= AuthorNameMaxLength, "name", fmt.Sprintf("must be less than %d bytes long", AuthorNameMaxLength))

	v.Check(strings.TrimSpace(a.Bio) != "", "bio", "must be provided")
	v.Check(len(a.Bio) <= AuthorNameMaxLength, "bio", fmt.Sprintf("must be less than %d bytes long", AuthorBioMaxLength))
}
