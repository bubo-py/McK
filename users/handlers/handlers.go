package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/service"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

var badRequestReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrBadRequest.ErrorType,
	ErrorMessage: customErrors.ErrBadRequest.Error(),
}

var notFoundReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrNotFound.ErrorType,
	ErrorMessage: customErrors.ErrNotFound.Error(),
}

var unauthenticatedReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrUnauthenticated.ErrorType,
	ErrorMessage: customErrors.ErrUnauthenticated.Error(),
}

var unauthorizedReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrUnauthorized.ErrorType,
	ErrorMessage: customErrors.ErrUnauthorized.Error(),
}

var unexpectedReturn = customErrors.ReturnError{
	ErrorType:    customErrors.ErrUnexpected.ErrorType,
	ErrorMessage: customErrors.ErrUnexpected.Error(),
}

type Handler struct {
	bl service.BusinessLogicInterface
}

func InitHandler(bl service.BusinessLogicInterface) Handler {
	var h Handler
	h.bl = bl
	return h
}

func (h Handler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u types.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	u, err = h.bl.AddUser(r.Context(), u)
	if err != nil {
		errBasedReturn(w, err)
		fmt.Println(err)
		return
	}

	err = json.NewEncoder(w).Encode(u.ID)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = h.bl.DeleteUser(r.Context(), id)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}

	var u types.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
		return
	}
	u.ID = id

	u, err = h.bl.UpdateUser(r.Context(), u, id)
	if err != nil {
		errBasedReturn(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(u.ID)
	if err != nil {
		log.Println(err)
	}
}

func errBasedReturn(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, customErrors.ErrBadRequest):
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(badRequestReturn)
		if err != nil {
			log.Println(err)
		}
	case errors.Is(err, customErrors.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(notFoundReturn)
		if err != nil {
			log.Println(err)
		}
	case errors.Is(err, customErrors.ErrUnauthenticated):
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(unauthenticatedReturn)
		if err != nil {
			log.Println(err)
		}
	case errors.Is(err, customErrors.ErrUnauthorized):
		w.WriteHeader(http.StatusForbidden)
		err = json.NewEncoder(w).Encode(unauthorizedReturn)
		if err != nil {
			log.Println(err)
		}
	default:
		err = json.NewEncoder(w).Encode(unexpectedReturn)
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			log.Println(err)
		}
	}
}
