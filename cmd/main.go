package main

import (
	"github.com/bubo-py/McK/events"
	"github.com/bubo-py/McK/handlers"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/urfave/cli/v2"
)

func main() {
	event1 := events.Event{
		ID:          1,
		Name:        "daily meeting",
		StartTime:   "2022-09-14T09:00:00.000Z",
		EndTime:     "2022-09-14T09:00:00.000Z",
		Description: "Friday daily meeting",
		AlertTime:   "2022-09-14T08:50:00.000Z",
	}

	events.AppendEvent(event1)
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

	r.Get("/api/events", handlers.GetEvents)
	r.Get("/api/events/{id}", handlers.GetEvent)
	r.Post("/api/events", handlers.AddEvent)
	r.Put("/api/events/{id}", handlers.UpdateEvent)
	r.Delete("/api/events/{id}", handlers.DeleteEvent)

	log.Fatal(http.ListenAndServe(":8080", r))
}
