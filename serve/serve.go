package serve

import (
	"context"
	"log"
	"net/http"
	"os"

	eventsHandlers "github.com/bubo-py/McK/events/handlers"
	eventsPostgres "github.com/bubo-py/McK/events/repositories/postgres"
	eventsService "github.com/bubo-py/McK/events/service"
	"github.com/bubo-py/McK/middlewares"
	usersHandlers "github.com/bubo-py/McK/users/handlers"
	usersPostgres "github.com/bubo-py/McK/users/repositories/postgres"
	usersService "github.com/bubo-py/McK/users/service"
	"github.com/go-chi/chi"
)

func Serve(ctx context.Context) {
	// Database setup
	connString := os.Getenv("PGURL")

	eventsDb, err := eventsPostgres.Init(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	usersDb, err := usersPostgres.Init(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	err = eventsPostgres.RunMigration(ctx, eventsDb)
	if err != nil {
		log.Fatal(err)
	}

	err = usersPostgres.RunMigration(ctx, usersDb)
	if err != nil {
		log.Fatal(err)
	}

	// Business logic setup
	eventsBl := eventsService.InitBusinessLogic(eventsDb)
	usersBl := usersService.InitBusinessLogic(usersDb)

	// Router setup
	r := chi.NewRouter()

	eventsHandler := eventsHandlers.InitHandler(eventsBl)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticate(usersBl))
		r.Mount("/api/events", eventsHandler.Mux)
	})

	usersHandler := usersHandlers.InitHandler(usersBl)
	r.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticate(usersBl))
		r.Mount("/api/users", usersHandler.Mux)
	})

	// Unprotected route
	r.Post("/api/users", usersHandler.AddUserHandler)

	port := os.Getenv("LISTEN_AND_SERVE_PORT")
	log.Printf("Starting an HTTP server on port %v", port)
	log.Fatal(http.ListenAndServe(port, r))

}
