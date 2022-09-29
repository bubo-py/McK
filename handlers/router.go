package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

var h Handler

func Serve() {
	r := chi.NewRouter()

	log.Println("Started an HTTP server on port 8080")

	r.Get("/api/events", h.GetEventsHandler)
	r.Get("/api/events/{id}", h.GetEventHandler)
	r.Post("/api/events", h.AddEventHandler)
	r.Put("/api/events/{id}", h.UpdateEventHandler)
	r.Delete("/api/events/{id}", h.DeleteEventHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
