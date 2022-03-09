package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrConflict       = errors.New("conflict")
	ErrURLDeleted     = errors.New("url was deleted")
)
