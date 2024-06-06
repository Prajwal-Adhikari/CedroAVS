package main

import "github.com/urfave/cli"


func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		{
			Name: "register-operator-with-eigenlayer",
			Aliases: []string{"regEl"},
			Usage: "Registers Verifier operator with EigenLayer",
			// Action: acti,
		},
	}
}