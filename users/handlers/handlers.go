package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/service"
	"github.com/go-chi/chi"
)

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
		log.Println(err)
	}

	u, err = h.bl.AddUser(r.Context(), u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
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
		log.Println(err)
	}

	err = h.bl.DeleteUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Println(err)
	}

	var u types.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println(err)
	}

	u, err = h.bl.UpdateUser(r.Context(), u, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(u.ID)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var u types.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Println(err)
	}

	err = h.bl.LoginUser(r.Context(), u.Login, u.Password)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode("logged in")
}
