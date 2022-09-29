package main

import (
	"log"
	"os"

	"github.com/bubo-py/McK/handlers"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{}

	app.Commands = []*cli.Command{
		{
			Name:  "serve",
			Usage: "start the HTTP service",
			Action: func(*cli.Context) error {
				handlers.Serve()
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
