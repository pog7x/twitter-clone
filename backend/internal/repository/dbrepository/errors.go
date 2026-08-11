package dbrepository

import "errors"

var (
	ErrNotFound     = errors.New("record not found")
	ErrInternal     = errors.New("internal repo error")
	ErrSelfFollow   = errors.New("user cannot follow itself")
	ErrInvalidMedia = errors.New("media not found, not owned by the user or already attached")
)
