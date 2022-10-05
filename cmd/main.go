package main

import (
	"context"
	"log"
	"os"

	"github.com/bubo-py/McK/handlers"
	"github.com/bubo-py/McK/repositories"
	"github.com/bubo-py/McK/service"
	"github.com/urfave/cli/v2"
)

func main() {
	db := repositories.InitDatabase()
	bl := service.InitBusinessLogic(db)
	handler := handlers.InitHandler(bl)

	ctx := context.Background()
	pg := repositories.PostgresInit(ctx)

	err := repositories.RunMigration(ctx, pg)
	if err != nil {
		log.Fatal(err)
	}

	//pg.AddEvent(ctx, types.Event{
	//	Name:        "the name",
	//	StartTime:   time.Time{},
	//	EndTime:     time.Time{},
	//	Description: "test",
	//	AlertTime:   time.Time{},
	//})

	app := &cli.App{}

	app.Commands = []*cli.Command{
		{
			Name:  "serve",
			Usage: "start the HTTP service",
			Action: func(*cli.Context) error {
				handlers.Serve(handler)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
