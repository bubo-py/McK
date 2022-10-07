package serve

import (
	"context"
	"log"

	"github.com/bubo-py/McK/handlers"
	"github.com/bubo-py/McK/repositories/postgres"
	"github.com/bubo-py/McK/service"
)

func Serve(ctx context.Context) {
	db, err := postgres.PostgresInit(ctx)
	if err != nil {
		log.Fatal(err)
	}

	bl := service.InitBusinessLogic(db)

	err = postgres.RunMigration(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.InitHandler(bl)
	handlers.InitRouter(handler)

}
