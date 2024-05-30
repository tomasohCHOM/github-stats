package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/urfave/cli/v2"
)

func main() {
	cmd := &cli.App{
		Name:  "greet",
		Usage: "say a greeting",
		Action: func(cli *cli.Context) error {
			client := github.NewClient(nil)
			response, _, err := client.Repositories.List(context.Background(), "tomasohCHOM", nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(response)
			return nil
		},
	}
	cmd.Run(os.Args)
}
