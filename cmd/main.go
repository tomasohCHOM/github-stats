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
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load .env file")
	}
	token := os.Getenv("ACCESS_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable not set")
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	owner := "EthanThatOneKid" // Replace with the repository owner
	repo := "acmcsuf.com"      // Replace with the repository name

	openIssuesCount, err := getIssueStats(ctx, client, owner, repo, "open")
	if err != nil {
		log.Fatalf("Error fetching open issues count: %v", err)
	}
	fmt.Printf("Open Issues: %d\n", openIssuesCount)

	closedIssuesCount, err := getIssueStats(ctx, client, owner, repo, "closed")
	if err != nil {
		log.Fatalf("Error fetching closed issues count: %v", err)
	}
	fmt.Printf("Closed Issues: %d\n", closedIssuesCount)

	openPRCount, err := getPRStats(ctx, client, owner, repo, "open")
	if err != nil {
		log.Fatalf("Error fetching open PRs count: %v", err)
	}
	fmt.Printf("Open PRs: %d\n", openPRCount)

	closedPRCount, err := getPRStats(ctx, client, owner, repo, "closed")
	if err != nil {
		log.Fatalf("Error fetching closed PRs count: %v", err)
	}
	fmt.Printf("Closed PRs: %d\n", closedPRCount)

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
	username := ctx.String("username")
	client := github.NewClient(nil)
	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 10}}

	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Repositories.List(context.Background(), username, opt)
		if err != nil {
			log.Fatal(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, repository := range allRepos {
		fmt.Println(repository.GetName())
	}

	return nil
}
