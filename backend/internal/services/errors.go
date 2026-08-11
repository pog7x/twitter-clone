package services

import "fmt"

type NotFoundServiceError struct {
	Err    error
	Entity string
}

func (e NotFoundServiceError) Error() string {
	return fmt.Sprintf("%s not found: %v", e.Entity, e.Err)
}

func (e NotFoundServiceError) Unwrap() error {
	return e.Err
}

type InternalServiceError struct {
	Err error
}

func (e InternalServiceError) Error() string {
	return fmt.Sprintf("internal error: %v", e.Err)
}

func (e InternalServiceError) Unwrap() error {
	return e.Err
}

type NotPermittedError struct{}

func (e NotPermittedError) Error() string {
	return "not permitted"
}

// InvalidCredentialsError is returned for any failed login attempt. It stays
// deliberately vague so the endpoint cannot be used to enumerate usernames.
type InvalidCredentialsError struct{}

func (e InvalidCredentialsError) Error() string {
	return "invalid username or password"
}
