package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
	"github.com/tomasohCHOM/github-stats/cmd/program"
	"github.com/tomasohCHOM/github-stats/cmd/ui/multiselector"
	"github.com/tomasohCHOM/github-stats/cmd/ui/selector"
	"github.com/tomasohCHOM/github-stats/cmd/ui/styles"
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

// Executes the CLI application
func ExecuteCLI(c *cli.Context) error {
	ctx := context.Background()

	var p *tea.Program

	userOptions := &program.ProgramState{
		Username:      c.String("username"),
		SelectedStats: []string{},
	}

	fmt.Printf("%s\n", styles.LogoStyle.Render(logo))

	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		err := fmt.Errorf("Please set your GitHub token in the ACCESS_TOKEN environment variable.")
		return err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	for {
		// Ask for the GitHub username to retrieve stats from
		if userOptions.Username == "" {
			header := "Which user would you like to retrieve GitHub data from?"
			errMsg := ""
			for {
				p = tea.NewProgram(text.InitialTextModel(userOptions, header, errMsg))
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
		p = tea.NewProgram(multiselector.InitialMultiSelectModel(userOptions, header, program.StatsOptions))
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		userOptions.ExitIfRequested(p)

		userOptions.RetrieveData(ctx, client)

		header = "Finished fetching the data. What would you like to do now?"
		p = tea.NewProgram(selector.InitialSelectionModel(userOptions, header, program.ContinueProgramOptions))
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return err
		}

		userOptions.ExecuteAfterRetrieval()
		if userOptions.ExitAfterContinueOption() {
			break
		}
	}

	return nil
}
