package entity

import "errors"

var (
	ErrInvalidEntity   = errors.New("invalid entity")
	ErrNotFound        = errors.New("not found")
	ErrNothingToUpdate = errors.New("nothing to update")
)
