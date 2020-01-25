package errors

type ValidateError struct {
	Message string
}

func NewValidateError(err string) ValidateError {
	return ValidateError{err}
}
