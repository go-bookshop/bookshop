package data

import "errors"

var (
	ErrDuplicateItem = errors.New("item must be unique")
)
