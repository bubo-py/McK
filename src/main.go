package main

import (
	"log"
	"net/http"
	"os"

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
