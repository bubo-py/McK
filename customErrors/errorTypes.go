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

var Unauthorized = CustomError{
	Err: errors.New("please provide your credentials"),
}

var ErrDB = CustomError{
	Err: errors.New("an error occurred in database"),
}
