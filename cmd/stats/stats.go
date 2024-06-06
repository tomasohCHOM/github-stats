package stats

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// GetRepositories fetches all the repository names associated with the GitHub username
func GetRepositories(ctx context.Context, client *github.Client, username string) ([]*github.Repository, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	totalRepos := []*github.Repository{}
	for {
		repos, response, err := client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}
		totalRepos = append(totalRepos, repos...)
		if response.NextPage == 0 {
			break
		}
		opts.Page = response.NextPage
	}
	return totalRepos, nil
}

// GetRepositoriesCount fetches the count of all repositories associated with the GitHub username
func GetRepositoriesCount(ctx context.Context, client *github.Client, username string) (int, error) {
	repos, err := GetRepositories(ctx, client, username)
	if err != nil {
		return 0, err
	}
	return len(repos), nil
}

// GetStarCount fetches the total number of stars earned by the associated GitHub user
func GetStarsCount(ctx context.Context, client *github.Client, username string) (int, error) {
	repos, err := GetRepositories(ctx, client, username)
	if err != nil {
		return 0, err
	}
	starCount := 0
	for _, repo := range repos {
		starCount += repo.GetStargazersCount()
	}
	return starCount, nil
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
