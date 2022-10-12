package serve

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bubo-py/McK/events/handlers"
	"github.com/bubo-py/McK/events/repositories/postgres"
	"github.com/bubo-py/McK/events/service"
	"github.com/bubo-py/McK/router"
)

func Serve(ctx context.Context) {
	connString := os.Getenv("PGURL")

	db, err := postgres.Init(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	bl := service.InitBusinessLogic(db)

	err = postgres.RunMigration(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.InitHandler(bl)
	r := router.InitRouter(handler)

	port := os.Getenv("LISTEN_AND_SERVE_PORT")
	log.Printf("Starting an HTTP server on port %v", port)
	log.Fatal(http.ListenAndServe(port, r))

}
