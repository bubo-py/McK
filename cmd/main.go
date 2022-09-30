package main

import (
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
