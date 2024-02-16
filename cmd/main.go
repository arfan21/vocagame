package main

import (
	"os"

	"github.com/arfan21/vocagame/cmd/api"
	migration "github.com/arfan21/vocagame/cmd/migrate"
	"github.com/urfave/cli/v2"
)

// @title Voca Game API
// @version 1.0
// @description This is a sample server cell for Voca Game Test API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.synapsis.id
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	appCli := cli.NewApp()
	appCli.Name = "Voca Game Test"
	appCli.Usage = "Voca Game Test API"
	appCli.Commands = []*cli.Command{
		migration.Root(),
		api.Serve(),
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}
