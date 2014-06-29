package main

import "github.com/dynport/dgtk/github"

func truncate(s string, l int, dots bool) string {
	if len(s) > l {
		if l > 6 && dots {
			return s[0:l-3] + "..."
		}
		return s[0:l]
	}
	return s
}

func loadIssue(client *github.Client, id int) (*github.Issue, error) {
	repo, e := githubRepo()
	if e != nil {
		return nil, e
	}
	a := github.LoadIssue{Repo: repo, Number: id}
	return a.Execute(client)
}

func loadIssues(client *github.Client, repo string) ([]*github.Issue, error) {
	a := &github.ListIssues{Repo: repo}
	return a.Execute(client)
}

// https://developer.github.com/v3/issues/#create-an-issue
