package routers

import (
	"net/http"

	"github.com/bubo-py/McK/events/handlers"
	"github.com/go-chi/chi"
)

func EventsRoutes(h handlers.Handler) http.Handler {
	r := chi.NewRouter()

	// events
	r.Get("/", h.GetEventsHandler)
	r.Get("/{id}", h.GetEventHandler)
	r.Post("/", h.AddEventHandler)
	r.Put("/{id}", h.UpdateEventHandler)
	r.Delete("/{id}", h.DeleteEventHandler)

	return r
}
