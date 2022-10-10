package handlers

import (
	"github.com/go-chi/chi"
)

func InitRouter(h Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/events", h.GetEventsHandler)
	r.Get("/api/events/{id}", h.GetEventHandler)
	r.Post("/api/events", h.AddEventHandler)
	r.Put("/api/events/{id}", h.UpdateEventHandler)
	r.Delete("/api/events/{id}", h.DeleteEventHandler)

	return r
}
