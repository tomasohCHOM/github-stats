// Run:
// go run cli/main.go --username tomasohCHOM

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"github.com/tomasohCHOM/github-stats"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load .env file")
	}
	app := &cli.App{
		Name:  "github-stats",
		Usage: "Generate GitHub stats from a user/organization",
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
	c := context.Background()

	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(c, ts)

	client := github.NewClient(tc)

	owner := ctx.String("username")

	opts := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 10}}

	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.List(context.Background(), owner, opts)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	for _, repository := range allRepos {
		repoName := repository.GetName()
		fmt.Println(repoName)
		openIssuesCount, err := githubstats.GetIssueStats(c, client, owner, repoName, "open")
		if err != nil {
			log.Fatalf("Error fetching open issues count: %v", err)
		}
		fmt.Printf("Open Issues: %d\n", openIssuesCount)

		closedIssuesCount, err := githubstats.GetIssueStats(c, client, owner, repoName, "closed")
		if err != nil {
			log.Fatalf("Error fetching closed issues count: %v", err)
		}
		fmt.Printf("Closed Issues: %d\n", closedIssuesCount)

		openPRCount, err := githubstats.GetPRStats(c, client, owner, repoName, "open")
		if err != nil {
			log.Fatalf("Error fetching open PRs count: %v", err)
		}
		fmt.Printf("Open PRs: %d\n", openPRCount)

		closedPRCount, err := githubstats.GetPRStats(c, client, owner, repoName, "closed")
		if err != nil {
			log.Fatalf("Error fetching closed PRs count: %v", err)
		}
		fmt.Printf("Closed PRs: %d\n", closedPRCount)

		fmt.Println("")
	}

	return nil
}
