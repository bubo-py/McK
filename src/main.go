package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/urfave/cli/v2"
)

var db []Event

func main() {
	event1 := Event{
		ID:          1,
		Name:        "daily meeting",
		StartTime:   "2022-09-14T09:00:00.000Z",
		EndTime:     "2022-09-14T09:00:00.000Z",
		Description: "Friday daily meeting",
		AlertTime:   "2022-09-14T08:50:00.000Z",
	}

	AppendEvent(event1)
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Get("/api/events", GetEvents)
	r.Post("/api/events", AddEvent)
	r.Put("/api/events/{id}", UpdateEvent)
	r.Delete("/api/events/{id}", DeleteEvent)

	log.Fatal(http.ListenAndServe(":8080", r))
}
