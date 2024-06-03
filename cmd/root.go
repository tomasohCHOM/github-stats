package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/tomasohCHOM/github-stats/cmd/ui/selector"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
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
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	ctx := context.Background()

	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Please set your GitHub token in the ACCESS_TOKEN environment variable.")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner := c.String("username")
	if owner == "" {
		fmt.Println("Which user/organization would you like to retrieve GitHub data from?")
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			owner = input
			_, _, err := client.Users.Get(ctx, owner)
			if err == nil {
				break
			}
			fmt.Println("User does not exist, try again!")
		}
	}

	header := "What would you like to retrieve?\n\n"
	options := []string{"Pull Request Count", "Issue Count"}

	p := bubbletea.NewProgram(selector.InitialSelectionModel(header, options))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}
	return nil
}
