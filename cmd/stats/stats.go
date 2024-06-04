package stats

import (
	"context"

	"github.com/google/go-github/github"
)

// GetRepositories fetches all the repository names and repository count associated with the GitHub username
func GetRepositories(ctx context.Context, client *github.Client, username string) ([]string, int, error) {
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository
	for {
		repo, resp, err := client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, 0, err
		}
		allRepos = append(allRepos, repo...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	repoNames := []string{}
	for _, repo := range allRepos {
		repoNames = append(repoNames, repo.GetName())
	}
	return repoNames, len(repoNames), nil
}

// GetIssueStats fetches the count of issues or pull requests based on the state and pull request filter
func GetIssueStats(ctx context.Context, client *github.Client, owner, repo, state string) (int, error) {
	opts := &github.IssueListByRepoOptions{
		State:       state,
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allIssues []*github.Issue
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
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

// GetPRStats fetches the count of issues or pull requests based on the state and pull request filter
func GetPRStats(ctx context.Context, client *github.Client, owner, repo, state string) (int, error) {
	opts := &github.PullRequestListOptions{
		State:       state,
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allPullRequests []*github.PullRequest
	for {
		pullRequests, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return 0, err
		}
		allPullRequests = append(allPullRequests, pullRequests...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return len(allPullRequests), nil
}
