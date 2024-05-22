package services

import "fmt"

type NotFoundServiceError struct {
	Err    error
	Entity string
}

func (e NotFoundServiceError) Error() string {
	return fmt.Sprintf("%#v: %s", e.Err, e.Entity)
}

type InternalServiceError struct {
	Err error
}

func (e InternalServiceError) Error() string {
	return fmt.Sprintf("internal error: %s", e.Err)
}

type NotPermittedError struct{}

func (e NotPermittedError) Error() string {
	return fmt.Sprintf("not permitted")
}
