package main

import (
	"log"
	"os"

	"github.com/joshpauline/vertex/internal/cli-tools"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Vertex CLI",
		Usage: "CLI tool to generate JSON format for Vertex GQL proxy",
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate vertex JSON format",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "schema",
						Aliases:     []string{"s"},
						Usage:       "Local schema file path",
						DefaultText: "schema.graphql",
						TakesFile:   true,
						Required:    true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "Generated output filepath",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "url",
						Aliases:  []string{"u"},
						Usage:    "Url of graphql service",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "Name of graphql service",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {

					schema := cCtx.Value("schema").(string)
					output := cCtx.Value("output").(string)
					url := cCtx.Value("url").(string)
					name := cCtx.Value("name").(string)

					err := cli_tools.GenerateVertex(schema, output, name, url)

					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}