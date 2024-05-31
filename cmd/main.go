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

	openIssuesCount, err := getStats(ctx, client, owner, repo, "open")
	if err != nil {
		log.Fatalf("Error fetching open issues count: %v", err)
	}
	fmt.Printf("Open Issues: %d\n", openIssuesCount)

	closedIssuesCount, err := getStats(ctx, client, owner, repo, "closed")
	if err != nil {
		log.Fatalf("Error fetching closed issues count: %v", err)
	}
	fmt.Printf("Closed Issues: %d\n", closedIssuesCount)

	// app := &cli.App{
	// Name:  "github-stats",
	// Usage: "Generate GitHub stats from a user/organization",
	//:! Flags: []cli.Flag{
	// &cli.StringFlag{
	// Name:     "username",
	// Usage:    "Match the github username",
	// Required: true,
	// },
	// },
	// Action: action,
	// }
	// if err := app.Run(os.Args); err != nil {
	// log.Fatal(err)
	// }
}

// getStats fetches the count of issues or pull requests based on the state and pull request filter
func getStats(ctx context.Context, client *github.Client, owner, repo, state string) (int, error) {
	opts := &github.IssueListByRepoOptions{
		State:       state,
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			log.Fatal(err)
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	count := 0
	for _, issue := range allIssues {
		if !issue.IsPullRequest() {
			count++
		}
	}

	return count, nil
}

func action(ctx *cli.Context) error {
	username := ctx.String("username")
	client := github.NewClient(nil)
	opt := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 10}}

	var allRepos []*github.Repository
	totalPRs := 0

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
		prs, _, err := client.PullRequests.List(context.Background(), username, repository.GetName(), &github.PullRequestListOptions{State: "open, closed"})
		if err != nil {
			log.Fatal(err)
		}
		totalPRs += len(prs)

		fmt.Println(repository.GetName())
		fmt.Println("Total Pull Requests: ", len(prs))
	}

	return nil
}
