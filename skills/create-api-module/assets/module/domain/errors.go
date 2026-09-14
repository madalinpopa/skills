package domain

import "errors"

var (
	ErrNotFound  = errors.New("{{ENTITY}} not found")
	ErrEmptyName = errors.New("{{ENTITY}} name is empty")
)
