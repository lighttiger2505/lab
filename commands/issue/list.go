package issue

import (
	"fmt"
	"strings"

	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

func makeIssueOption(issueListOption *ListIssueOption) *gitlab.ListIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listIssuesOptions := &gitlab.ListIssuesOptions{
		State:       gitlab.String(issueListOption.GetState()),
		Scope:       gitlab.String(issueListOption.GetScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listIssuesOptions
}

func makeProjectIssueOption(issueListOption *ListIssueOption) *gitlab.ListProjectIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(issueListOption.GetState()),
		Scope:       gitlab.String(issueListOption.GetScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func listOfProject(client lab.Issue, project string, opt *ListIssueOption) (string, error) {
	issues, err := client.GetProjectIssues(
		makeProjectIssueOption(opt),
		project,
	)
	if err != nil {
		return "", err
	}

	// Print issue list
	output := projectIssueOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

func listAll(client lab.Issue, opt *ListIssueOption) (string, error) {
	issues, err := client.GetAllProjectIssues(makeIssueOption(opt))
	if err != nil {
		return "", err
	}

	// Print issue list
	output := issueOutput(issues)
	result := columnize.SimpleFormat(output)
	return result, nil
}

func issueOutput(issues []*gitlab.Issue) []string {
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

func projectIssueOutput(issues []*gitlab.Issue) []string {
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
