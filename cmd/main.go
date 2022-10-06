package main

import (
	"context"
	"log"
	"os"

	"github.com/bubo-py/McK/serve"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx := context.Background()

	app := &cli.App{}

	app.Commands = []*cli.Command{
		{
			Name:  "serve",
			Usage: "start the HTTP service",
			Action: func(*cli.Context) error {
				serve.Serve(ctx)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
