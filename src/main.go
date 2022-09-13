package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	serve()
}

func serve() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
