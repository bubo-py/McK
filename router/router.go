package router

import (
	"github.com/bubo-py/McK/events/handlers"
	"github.com/go-chi/chi"
)

func InitRouter(h handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/events", h.GetEventsHandler)
	r.Get("/api/events/{id}", h.GetEventHandler)
	r.Post("/api/events", h.AddEventHandler)
	r.Put("/api/events/{id}", h.UpdateEventHandler)
	r.Delete("/api/events/{id}", h.DeleteEventHandler)

	return r
}
