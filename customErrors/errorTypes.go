package customErrors

import (
	"errors"
)

type CustomError struct {
	Err       error
	ErrorType string
}

func (ce CustomError) Error() string {
	return ce.Err.Error()
}

var BadRequest = CustomError{
	Err:       errors.New("the server cannot process the request"),
	ErrorType: "BadRequest",
}

var IncorrectCredentials = CustomError{
	Err: errors.New("incorrect credentials"),
}

var ErrUnauthorized = CustomError{
	Err: errors.New("please provide your credentials"),
}

var ErrUnexpected = CustomError{
	Err:       errors.New("an unexpected error occurred"),
	ErrorType: "Unexpected",
}
