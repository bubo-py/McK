package serve

import (
	"context"
	"log"
	"net/http"
	"os"

	eventsHandlers "github.com/bubo-py/McK/events/handlers"
	eventsPostgres "github.com/bubo-py/McK/events/repositories/postgres"
	eventsRouters "github.com/bubo-py/McK/events/routers"
	eventsService "github.com/bubo-py/McK/events/service"
	usersHandlers "github.com/bubo-py/McK/users/handlers"
	usersPostgres "github.com/bubo-py/McK/users/repositories/postgres"
	userRouters "github.com/bubo-py/McK/users/routers"
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
	r.Mount("/api/events", eventsRouters.EventsRoutes(eventsHandler))

	usersHandler := usersHandlers.InitHandler(usersBl)
	r.Mount("/api/users", userRouters.UserRoutes(usersHandler))

	port := os.Getenv("LISTEN_AND_SERVE_PORT")
	log.Printf("Starting an HTTP server on port %v", port)
	log.Fatal(http.ListenAndServe(port, r))

}
