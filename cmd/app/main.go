package main

import (
	"os"

	"github.com/optclblast/currencier/internal/app"
	"github.com/optclblast/currencier/internal/config"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "currencier",
		Version: "0.1.0",
		Flags: []cli.Flag{
			// common
			&cli.StringFlag{
				Name:  "log-level",
				Value: "debug",
			},

			// rest
			&cli.StringFlag{
				Name:  "http-port",
				Value: "8080",
			},

			// database
			&cli.StringFlag{
				Name:  "cache-url",
				Value: "localhost:6379",
			},
			&cli.StringFlag{
				Name: "cache-user",
			},
			&cli.StringFlag{
				Name: "cache-secret",
			},
		},
		Action: func(c *cli.Context) error {
			app.Run(&config.Config{
				Common: config.CommonConfig{
					Level: c.String("log-level"),
				},
				Rest: config.RestConfig{
					Port: c.String("http-port"),
				},
				Cache: config.Cache{
					URL:    c.String("cache-url"),
					User:   c.String("cache-user"),
					Secret: c.String("cache-secret"),
				},
			})

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
