package serve

import (
	"context"
	"log"

	"github.com/bubo-py/McK/handlers"
	"github.com/bubo-py/McK/repositories"
	"github.com/bubo-py/McK/service"
)

func Serve(ctx context.Context) {
	db := repositories.PostgresInit(ctx)
	bl := service.InitBusinessLogic(db)

	handler := handlers.InitHandler(bl)
	handlers.InitRouter(handler)

	err := repositories.RunMigration(ctx, db)
	if err != nil {
		log.Fatal(err)
	}
}
