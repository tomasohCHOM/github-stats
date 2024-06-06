package program

import (
	"context"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/github"
	"github.com/tomasohCHOM/github-stats/cmd/stats"
)

var statsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7CE38B")).Bold(true)

const (
	STARS        = "Earned Stars Count"
	PRS          = "Pull Requests Count"
	ISSUES       = "Issues Count"
	REPOSITORIES = "Repositories Count"
)

var RetrievalOptions = []string{STARS, PRS, ISSUES, REPOSITORIES}

type ProgramState struct {
	ExitState     bool
	Username      string
	DataToCollect []string
}

func (p *ProgramState) ExitIfRequested(tprogram *tea.Program) {
	if p.ExitState {
		if err := tprogram.ReleaseTerminal(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Good Bye!")
		os.Exit(1)
	}
}

func (p *ProgramState) RetrieveData(ctx context.Context, client *github.Client) {
	for _, selection := range p.DataToCollect {
		var data int
		var err error

		switch selection {
		case STARS:
			data, err = stats.GetStarsCount(ctx, client, p.Username)
		case PRS:
			data, err = stats.GetPRStats(ctx, client, p.Username)
		case ISSUES:
			data, err = stats.GetIssueStats(ctx, client, p.Username)
		case REPOSITORIES:
			data, err = stats.GetRepositoriesCount(ctx, client, p.Username)
		}

		if err != nil {
			log.Fatal("Could not fetch data")
		}
		fmt.Printf("%s: %s\n", selection, statsStyle.Render(fmt.Sprintf("%d", data)))
	}
}
