package customErrors

import (
	"errors"
)

type ReturnError struct {
	ErrorType    string
	ErrorMessage string
}

type CustomError struct {
	Err       error
	ErrorType string
}

func (ce CustomError) Error() string {
	return ce.Err.Error()
}

var ErrBadRequest = CustomError{
	Err:       errors.New("the server cannot process the request"),
	ErrorType: "BadRequest",
}

var ErrNotFound = CustomError{
	Err:       errors.New("the server cannot find the requested resource"),
	ErrorType: "NotFound",
}

var ErrUnauthenticated = CustomError{
	Err: errors.New("failed to authenticate current user"),
	// Authentication confirms that users are who they say they are, 401
	ErrorType: "Unauthenticated",
}

var ErrUnauthorized = CustomError{
	Err: errors.New("the server cannot process the request due to lack of client's access rights"),
	// Authorization checks whether users have permission to access a resource, 403
	ErrorType: "Unauthorized",
}

var ErrUnexpected = CustomError{
	Err:       errors.New("an unexpected error occurred"),
	ErrorType: "Unexpected",
}
