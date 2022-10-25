package customErrors

import (
	"errors"
)

type CustomError struct {
	Err error
}

func (ce CustomError) Error() string {
	return ce.Err.Error()
}

var BadRequest = CustomError{
	Err: errors.New("the server cannot process the request"),
}

var IncorrectCredentials = CustomError{
	Err: errors.New("incorrect credentials"),
}

var ErrUnauthorized = CustomError{
	Err: errors.New("please provide your credentials"),
}

var ErrUnexpected = CustomError{
	Err: errors.New("an unexpected error occurred"),
}
