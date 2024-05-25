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
				Name: "http-port",
			},

			// database
			&cli.StringFlag{
				Name: "db-url",
			},
			&cli.IntFlag{
				Name:  "pool",
				Value: 2,
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
				PG: config.PG{
					URL:     c.String("db-url"),
					PoolMax: c.Int("pool"),
				},
			})

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
