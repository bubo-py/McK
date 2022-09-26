package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bubo-py/McK/handlers"
	"github.com/go-chi/chi"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{}

	app.Commands = []*cli.Command{
		{
			Name:  "serve",
			Usage: "start the HTTP service",
			Action: func(*cli.Context) error {
				serve()
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func serve() {
	r := chi.NewRouter()

	log.Println("Started an HTTP server on port 8080")

	r.Get("/api/events", handlers.GetEventsHandler)
	r.Get("/api/events/{id}", handlers.GetEventHandler)
	r.Post("/api/events", handlers.AddEventHandler)
	r.Put("/api/events/{id}", handlers.UpdateEventHandler)
	r.Delete("/api/events/{id}", handlers.DeleteEventHandler)

	r.Get("/api/events/filters/day/{day}", handlers.GetEventsByDay)
	r.Get("/api/events/filters/month/{month}", handlers.GetEventsByMonth)
	r.Get("/api/events/filters/year/{year}", handlers.GetEventsByYear)

	log.Fatal(http.ListenAndServe(":8080", r))
}
