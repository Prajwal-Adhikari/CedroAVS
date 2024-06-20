package main

import (
	"log"
	"os"

	"github.com/Cedro-Finance/CedroAVS/verifier/cli/actions"
	"github.com/urfave/cli"
)

var config string

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{cli.StringFlag{
		Name:     "config",
		Required: true,
		Usage:    "Load configuration from `FILE`",
	}}

	app.Commands = []cli.Command{
		{
			Name:    "register-operator-with-avs",
			Aliases: []string{"reg"},
			Usage:   "Registers Verifier operator with EigenLayer",
			Action:  actions.RegisterOperatorWithAvs,
		},
		{
			Name:    "deregister-operator-with-avs",
			Aliases: []string{"dereg"},
			Usage:   "Registers Verifier operator with EigenLayer",
			Action:  actions.DeregisterOperatorWithAvs,
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
