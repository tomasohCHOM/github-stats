package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/github"
	"github.com/tomasohCHOM/github-stats/cmd/program"
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

var (
	logoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#77BDFB")).Bold(true)
	headerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A2D2FB")).Bold(true)
	contrastStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CEA5FB")).Bold(true)
)

func ExecuteCLI(c *cli.Context) error {
	ctx := context.Background()

	userOptions := &program.ProgramState{
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
			p := tea.NewProgram(text.InitialTextModel(userOptions, header, errMsg))
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				return err
			}

			userOptions.ExitIfRequested(p)

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

	p := tea.NewProgram(selector.InitialSelectionModel(userOptions, header, program.RetrievalOptions))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	fmt.Println(headerStyle.Render("\nFetching the data for you..."))
	userOptions.RetrieveData(ctx, client)

	return nil
}
