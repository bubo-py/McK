package routers

import (
	"net/http"

	"github.com/bubo-py/McK/users/handlers"
	"github.com/go-chi/chi"
)

func UserRoutes(h handlers.Handler) http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.AddUserHandler)
	r.Put("/{id}", h.UpdateUserHandler)
	r.Delete("/{id}", h.DeleteUserHandler)
	r.Post("/login", h.LoginHandler)

	return r
}
