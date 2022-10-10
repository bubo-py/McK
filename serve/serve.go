package serve

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/bubo-py/McK/handlers"
	"github.com/bubo-py/McK/repositories/postgres"
	"github.com/bubo-py/McK/service"
)

func Serve(ctx context.Context) {
	connString := os.Getenv("PGURL")

	db, err := postgres.PostgresInit(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	bl := service.InitBusinessLogic(db)

	err = postgres.RunMigration(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.InitHandler(bl)
	r := handlers.InitRouter(handler)

	port := os.Getenv("LISTEN_AND_SERVE_PORT")
	log.Printf("Starting an HTTP server on port %v", port)
	log.Fatal(http.ListenAndServe(port, r))

}
