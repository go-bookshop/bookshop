package data

import "errors"

var (
	ErrDuplicateAuthorName = errors.New("authors must be unique")
)
