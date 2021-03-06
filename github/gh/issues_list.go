package main

import (
	"fmt"
	"strings"

	"github.com/dynport/dgtk/github"
	"github.com/dynport/gocli"
)

type issuesList struct {
	Query     []string `cli:"arg"`
	All       bool     `cli:"opt --all"`
	Closed    bool     `cli:"opt --closed"`
	Assignee  string   `cli:"opt --assignee"`
	Creator   string   `cli:"opt --creator"`
	Mentioned string   `cli:"opt --mentioned"`
	Asc       bool     `cli:"opt --asc"`
	Sort      string   `cli:"opt --sort"`
	Labels    string   `cli:"opt --labels"`
	Milestone int      `cli:"opt --milestone"`
}

func (r *issuesList) Run() error {
	var e error
	repo, e := githubRepo()
	if e != nil {
		return e
	}
	a := &github.ListIssues{Assignee: r.Assignee, Creator: r.Creator, Mentioned: r.Mentioned, Sort: r.Sort, Repo: repo, Milestone: r.Milestone}
	if r.Labels != "" {
		a.Labels = strings.Split(r.Labels, ",")
	}
	if r.Asc {
		a.Direction = github.IssueSortAsc
	}
	if r.All {
		a.State = github.IssueStateAll
	} else if r.Closed {
		a.State = github.IssueStateClosed
	}

	cl, e := client()
	if e != nil {
		return e
	}

	issues, e := a.Execute(cl)
	if e != nil {
		return e
	}
	if len(issues) == 0 {
		fmt.Println("no issues found")
		return nil
	}
	t := gocli.NewTable()
	for _, i := range issues {
		matches := func() bool {
			for _, q := range r.Query {
				if !strings.Contains(i.Title, q) {
					return false
				}
			}
			return true
		}()
		if !matches {
			continue
		}
		orga, issueRepo, e := i.Repo()
		if e != nil {
			return e
		}
		labels := []string{}
		for _, l := range i.Labels {
			labels = append(labels, l.Name)
		}
		parts := []interface{}{i.Number}
		if repo == "" {
			parts = append(parts, orga+"/"+issueRepo)
		}
		assignee := ""
		if i.Assignee != nil {
			assignee = i.Assignee.Login
		}
		parts = append(parts, truncate(i.Title, 48, true), truncate(i.CreatedAt, 16, false), assignee, i.State, strings.Join(labels, ","))
		t.Add(parts...)
	}
	fmt.Println(t)
	return nil
}
