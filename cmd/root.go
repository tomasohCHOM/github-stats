package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/tomasohCHOM/github-stats/cmd/state"
	"github.com/tomasohCHOM/github-stats/cmd/stats"
	"github.com/tomasohCHOM/github-stats/cmd/ui/selector"
	"github.com/tomasohCHOM/github-stats/cmd/ui/text"

	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

const logo = `
   _____ _ _   _    _       _        _____ _        _       
  / ____(_) | | |  | |     | |      / ____| |      | |      
 | |  __ _| |_| |__| |_   _| |__   | (___ | |_ __ _| |_ ___ 
 | | |_ | | __|  __  | | | | '_ \   \___ \| __/ _, | __/ __|
 | |__| | | |_| |  | | |_| | |_) |  ____) | || (_| | |_\__ \
  \_____|_|\__|_|  |_|\__,_|_.__/  |_____/ \__\__,_|\__|___/
`

var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#77BDFB")).Bold(true)

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

	fmt.Printf("%s\n", logoStyle.Render(logo))

	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("Please set your GitHub token in the ACCESS_TOKEN environment variable.")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	if userOptions.Username == "" {
		header := "Which user/organization would you like to retrieve GitHub data from?"
		errMsg := ""
		for {
			p := bubbletea.NewProgram(text.InitialTextModel(userOptions, header, errMsg))
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}
			if userOptions.Username == "" {
				errMsg = "Username should not be empty, please try again."
				continue
			}
			username := userOptions.Username
			_, _, err := client.Users.Get(ctx, username)
			if err == nil {
				break
			}
			errMsg = "User does not exist, please try again."
		}
	}

	header := fmt.Sprintf("Analyzing username %s, what would you like to retrieve?", userOptions.Username)
	options := []string{"Repository Count", "Pull Request Count", "Issue Count"}

	p := bubbletea.NewProgram(selector.InitialSelectionModel(userOptions, header, options))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Println("\nWaiting for Response...")

	for _, selection := range userOptions.DataToCollect {
		switch selection {
		case "Repository Count":
			repoCount, err := stats.GetRepositoriesCount(ctx, client, userOptions.Username)
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			fmt.Println("Public Repository Count:", repoCount)

		case "Pull Request Count":
			prCount, err := stats.GetPRStats(ctx, client, userOptions.Username)
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			fmt.Println("Total PR Count:", prCount)

		case "Issue Count":
			issueCount, err := stats.GetIssueStats(ctx, client, userOptions.Username)
			if err != nil {
				log.Fatal("Could not fetch data")
			}
			fmt.Println("Total Issue Count:", issueCount)
		}
	}

	return nil
}
