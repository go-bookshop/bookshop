package repository

import "errors"

var (
	ErrDuplicateItem      = errors.New("item must be unique")
	ErrRecordNotFound     = errors.New("record not found")
	ErrRecordEditConflict = errors.New("record edit conflict")
)
