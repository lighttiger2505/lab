package mr

import (
	"strings"

	"github.com/fatih/color"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type listMethod struct {
	internal.Method
	client  lab.MergeRequest
	opt     *ListOption
	project string
}

func (m *listMethod) Process() (string, error) {
	mergeRequests, err := m.client.GetProjectMargeRequest(
		makeProjectMergeRequestOption(m.opt),
		m.project,
	)
	if err != nil {
		return "", nil
	}
	outputs := outProjectMergeRequest(mergeRequests)
	return columnize.SimpleFormat(outputs), nil
}

type listAllMethod struct {
	internal.Method
	client lab.MergeRequest
	opt    *ListOption
}

func (m *listAllMethod) Process() (string, error) {
	// Do get merge request list
	mergeRequests, err := m.client.GetAllProjectMergeRequest(
		makeMergeRequestOption(m.opt),
	)
	if err != nil {
		return "", nil
	}

	// Print merge request list
	outputs := outMergeRequest(mergeRequests)
	return columnize.SimpleFormat(outputs), nil
}

func makeMergeRequestOption(listMergeRequestsOption *ListOption) *gitlab.ListMergeRequestsOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: listMergeRequestsOption.Num,
	}

	listRequestsOptions := &gitlab.ListMergeRequestsOptions{
		State:       gitlab.String(listMergeRequestsOption.getState()),
		Scope:       gitlab.String(listMergeRequestsOption.getScope()),
		OrderBy:     gitlab.String(listMergeRequestsOption.OrderBy),
		Sort:        gitlab.String(listMergeRequestsOption.Sort),
		Search:      gitlab.String(listMergeRequestsOption.Search),
		ListOptions: *listOption,
	}

	if listMergeRequestsOption.Milestone != "" {
		listRequestsOptions.Milestone = gitlab.String(listMergeRequestsOption.Milestone)
	}
	if listMergeRequestsOption.AuthorID != 0 {
		listRequestsOptions.AuthorID = gitlab.Int(listMergeRequestsOption.AuthorID)
	}
	if listMergeRequestsOption.AssigneeID != 0 {
		listRequestsOptions.AssigneeID = gitlab.Int(listMergeRequestsOption.AssigneeID)
	}
	return listRequestsOptions
}

func makeProjectMergeRequestOption(listMergeRequestsOption *ListOption) *gitlab.ListProjectMergeRequestsOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: listMergeRequestsOption.Num,
	}
	listMergeRequestsOptions := &gitlab.ListProjectMergeRequestsOptions{
		State:       gitlab.String(listMergeRequestsOption.getState()),
		Scope:       gitlab.String(listMergeRequestsOption.getScope()),
		OrderBy:     gitlab.String(listMergeRequestsOption.OrderBy),
		Sort:        gitlab.String(listMergeRequestsOption.Sort),
		Search:      gitlab.String(listMergeRequestsOption.Search),
		ListOptions: *listOption,
	}

	if listMergeRequestsOption.Milestone != "" {
		listMergeRequestsOptions.Milestone = gitlab.String(listMergeRequestsOption.Milestone)
	}
	if listMergeRequestsOption.AuthorID != 0 {
		listMergeRequestsOptions.AuthorID = gitlab.Int(listMergeRequestsOption.AuthorID)
	}
	if listMergeRequestsOption.AssigneeID != 0 {
		listMergeRequestsOptions.AssigneeID = gitlab.Int(listMergeRequestsOption.AssigneeID)
	}
	return listMergeRequestsOptions
}

func outProjectMergeRequest(mergeRequsets []*gitlab.MergeRequest) []string {
	yellow := color.New(color.FgYellow).SprintFunc()
	outputs := []string{}
	for _, mergeRequest := range mergeRequsets {
		output := strings.Join([]string{
			yellow(mergeRequest.IID),
			mergeRequest.Title,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}

func outMergeRequest(mergeRequsets []*gitlab.MergeRequest) []string {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	outputs := []string{}
	for _, mergeRequest := range mergeRequsets {
		output := strings.Join([]string{
			cyan(internal.ParceRepositoryFullName(mergeRequest.WebURL)),
			yellow(mergeRequest.IID),
			mergeRequest.Title,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
