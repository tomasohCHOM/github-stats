// Run:
// go run cli/main.go --username tomasohCHOM

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
	app := &cli.App{
		Name:  "github-stats",
		Usage: "Generate GitHub stats from a user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Usage:    "Match the github username",
				Required: true,
			},
		},
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	username := ctx.String("username")
	client := github.NewClient(nil)
	response, _, err := client.Repositories.List(context.Background(), username, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
	return nil
}
