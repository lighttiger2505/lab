package issue

import (
	"strings"

	"github.com/fatih/color"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type listMethod struct {
	client  lab.Issue
	opt     *ListOption
	project string
}

func (m *listMethod) Process() (string, error) {
	issues, err := m.client.GetProjectIssues(
		makeProjectIssueOption(m.opt),
		m.project,
	)
	if err != nil {
		return "", err
	}

	output := listOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

type listAllMethod struct {
	client lab.Issue
	opt    *ListOption
}

func (m *listAllMethod) Process() (string, error) {
	issues, err := m.client.GetAllProjectIssues(makeAllProjectIssueOption(m.opt))
	if err != nil {
		return "", err
	}

	output := listAllOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

func makeProjectIssueOption(issueListOption *ListOption) *gitlab.ListProjectIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(issueListOption.getState()),
		Scope:       gitlab.String(issueListOption.getScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		Search:      gitlab.String(issueListOption.Search),
		Milestone:   gitlab.String(issueListOption.Milestone),
		ListOptions: *listOption,
	}
	if issueListOption.AuthorID != 0 {
		listProjectIssuesOptions.AuthorID = gitlab.Int(issueListOption.AuthorID)
	}
	return listProjectIssuesOptions
}

func makeAllProjectIssueOption(issueListOption *ListOption) *gitlab.ListIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listIssuesOptions := &gitlab.ListIssuesOptions{
		State:       gitlab.String(issueListOption.getState()),
		Scope:       gitlab.String(issueListOption.getScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		Search:      gitlab.String(issueListOption.Search),
		Milestone:   gitlab.String(issueListOption.Milestone),
		ListOptions: *listOption,
	}
	if issueListOption.AuthorID != 0 {
		listIssuesOptions.AuthorID = gitlab.Int(issueListOption.AuthorID)
	}
	return listIssuesOptions
}

func listOutput(issues []*gitlab.Issue) []string {
	yellow := color.New(color.FgYellow).SprintFunc()
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			yellow(issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}

func listAllOutput(issues []*gitlab.Issue) []string {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			cyan(internal.ParceRepositoryFullName(issue.WebURL)),
			yellow(issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}
