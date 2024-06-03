package stats

import (
	"context"
	"log"

	"github.com/google/go-github/github"
)

// getIssueStats fetches the count of issues or pull requests based on the state and pull request filter
func GetIssueStats(ctx context.Context, client *github.Client, owner, repo, state string) (int, error) {
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

// getPRStats fetches the count of issues or pull requests based on the state and pull request filter
func GetPRStats(ctx context.Context, client *github.Client, owner, repo, state string) (int, error) {
	opts := &github.PullRequestListOptions{
		State:       state,
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allPullRequests []*github.PullRequest
	for {
		pullRequests, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			log.Fatal(err)
		}
		allPullRequests = append(allPullRequests, pullRequests...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return len(allPullRequests), nil
}
