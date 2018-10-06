package issue

import (
	"fmt"
	"strings"

	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

func makeIssueOption(issueListOption *ListOption) *gitlab.ListIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listIssuesOptions := &gitlab.ListIssuesOptions{
		State:       gitlab.String(issueListOption.getState()),
		Scope:       gitlab.String(issueListOption.getScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listIssuesOptions
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
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func list(client lab.Issue, project string, opt *ListOption) (string, error) {
	issues, err := client.GetProjectIssues(
		makeProjectIssueOption(opt),
		project,
	)
	if err != nil {
		return "", err
	}

	// Print issue list
	output := listAllOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

func listAll(client lab.Issue, opt *ListOption) (string, error) {
	issues, err := client.GetAllProjectIssues(makeIssueOption(opt))
	if err != nil {
		return "", err
	}

	// Print issue list
	output := listOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

func listOutput(issues []*gitlab.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			lab.ParceRepositoryFullName(issue.WebURL),
			fmt.Sprintf("%d", issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}

func listAllOutput(issues []*gitlab.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("%d", issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}
