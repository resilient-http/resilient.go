package middlewares

import "fmt"

type MiddlewareError struct {
	code    int
	message string
	name    string
}

func NewMiddlewareError(name, message string, code int) error {
	return MiddlewareError{
		name:    name,
		message: message,
		code:    code,
	}
}

func (e MiddlewareError) Error() string {
	return fmt.Sprintf("Middleware error [%s]: %s", e.name, e.message)
}
