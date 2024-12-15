package models

import "time"

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
