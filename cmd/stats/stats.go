package stats

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// GetRepositories fetches all the repository names and repository count associated with the GitHub username
func GetRepositories(ctx context.Context, client *github.Client, username string) (int, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	totalRepos := 0
	for {
		repos, response, err := client.Repositories.List(ctx, username, opts)
		if err != nil {
			return 0, err
		}
		totalRepos += len(repos)
		if response.NextPage == 0 {
			break
		}
		opts.Page = response.NextPage
	}
	return totalRepos, nil
}

// GetIssueStats fetches the total issue count for the associated GitHub username
func GetIssueStats(ctx context.Context, client *github.Client, username string) (int, error) {
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	query := fmt.Sprintf("type:issue author:%s", username)
	issues, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		return 0, err
	}
	return *issues.Total, nil
}

// GetPRStats fetches the total pull request count for the associated GitHub username
func GetPRStats(ctx context.Context, client *github.Client, username string) (int, error) {
	opts := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	query := fmt.Sprintf("type:pr author:%s", username)
	prs, _, err := client.Search.Issues(ctx, query, opts)
	if err != nil {
		return 0, err
	}
	return *prs.Total, nil
}
