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
	"github.com/tomasohCHOM/github-stats/cmd/state"
	"github.com/tomasohCHOM/github-stats/cmd/stats"
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

	userOptions := &state.UserOptions{
		Username:      c.String("username"),
		DataToCollect: []string{},
	}

	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Please set your GitHub token in the ACCESS_TOKEN environment variable.")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if userOptions.Username == "" {
		fmt.Println("Which user/organization would you like to retrieve GitHub data from?")
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			username := input
			_, _, err := client.Users.Get(ctx, username)
			if err == nil {
				userOptions.Username = username
				break
			}
			fmt.Println("User does not exist, try again!")
		}
	}

	header := "What would you like to retrieve?\n\n"
	options := []string{"Repository Count", "Pull Request Count", "Issue Count"}

	p := bubbletea.NewProgram(selector.InitialSelectionModel(userOptions, header, options))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Println("\nWaiting for Response...")

	for _, selection := range userOptions.DataToCollect {
		repo := "lc-dailies"
		switch selection {
		case "Pull Request Count":
			prCount := 0
			open, err := stats.GetPRStats(ctx, client, userOptions.Username, repo, "open")
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			prCount += open
			closed, err := stats.GetPRStats(ctx, client, userOptions.Username, repo, "closed")
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			prCount += closed
			fmt.Println("Total PR Count (Open and Closed):", prCount)
		case "Issue Count":
			issueCount := 0
			open, err := stats.GetIssueStats(ctx, client, userOptions.Username, repo, "open")
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			issueCount += open
			closed, err := stats.GetIssueStats(ctx, client, userOptions.Username, repo, "closed")
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			issueCount += closed
			fmt.Println("Total PR Count (Open and Closed):", issueCount)
		}
	}

	return nil
}
