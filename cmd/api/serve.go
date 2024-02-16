package api

import (
	"github.com/arfan21/vocagame/config"
	"github.com/arfan21/vocagame/internal/server"
	dbpostgres "github.com/arfan21/vocagame/pkg/db/postgres"
	dbredis "github.com/arfan21/vocagame/pkg/db/redis"
	"github.com/urfave/cli/v2"
)

func Serve() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Run the API server",
		Action: func(c *cli.Context) error {
			_, err := config.LoadConfig()
			if err != nil {
				return err
			}

			_, err = config.ParseConfig(config.GetViper())
			if err != nil {
				return err
			}

			db, err := dbpostgres.NewPgx()
			if err != nil {
				return err
			}

			dbRedis, err := dbredis.New()
			if err != nil {
				return err
			}

			server := server.New(
				db,
				dbRedis,
			)
			return server.Run()
		},
	}

}
