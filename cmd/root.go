package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/urfave/cli/v2"
)

func Execute() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load .env file")
	}
	app := &cli.App{
		Name:  "github-stats",
		Usage: "Generate GitHub stats from a user/organization",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "username",
				Value: "",
				Usage: "Match the github username",
			},
		},
		Action: ExecuteCLI,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
