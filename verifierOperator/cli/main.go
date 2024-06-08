package main

import (
	"log"
	"os"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/cli/actions"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{cli.StringFlag{
		Name:     "config",
		Required: false,
		Usage:    "Load configuration from `FILE`",
	}}

	app.Commands = []cli.Command{
		{
			Name:    "register-operator-with-eigenlayer",
			Aliases: []string{"regEl"},
			Usage:   "Registers Verifier operator with EigenLayer",
			Action:  actions.RegisterOperatorWithAvs,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
