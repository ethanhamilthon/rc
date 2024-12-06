package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "hc",
		Usage:   "TUI based HTTP Client",
		Version: "v0.1.0",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Creates new HC project",
				Action: func(c *cli.Context) error {
					// arg := c.Args().Slice()
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			HandleRequest(c.Args().Slice())
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
