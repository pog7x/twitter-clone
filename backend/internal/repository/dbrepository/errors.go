package dbrepository

import "errors"

var ErrNotFound = errors.New("record not found")
var ErrInternal = errors.New("internal repo error")
