package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	cmd := &cli.App{
		Name:  "greet",
		Usage: "say a greeting",
		Action: func(cli *cli.Context) error {
			fmt.Println("Greetings")
			return nil
		},
	}
	cmd.Run(os.Args)
}
