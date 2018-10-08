package mr

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

type detailMethod struct {
	internal.Method
	client  lab.MergeRequest
	project string
	id      int
}

func (m *detailMethod) Process() (string, error) {
	// Do get merge request
	mergeRequest, err := m.client.GetMergeRequest(m.id, m.project)
	if err != nil {
		return "", err
	}
	return outMergeRequestDetail(mergeRequest), nil
}

func outMergeRequestDetail(mergeRequest *gitlab.MergeRequest) string {
	base := `%s %s [%s] (created by @%s, %s)
Assignee: %s
Milestone: %s
Labels: %s

%s`

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	var stateColor func(a ...interface{}) string
	if mergeRequest.State == "opened" {
		stateColor = color.New(color.FgGreen).SprintFunc()
	} else {
		stateColor = color.New(color.FgRed).SprintFunc()
	}

	milestone := ""
	if mergeRequest.Milestone != nil {
		milestone = mergeRequest.Milestone.Title
	}

	detial := fmt.Sprintf(base,
		yellow(mergeRequest.IID),
		cyan(mergeRequest.Title),
		stateColor(mergeRequest.State),
		mergeRequest.Author.Name,
		mergeRequest.CreatedAt.String(),
		mergeRequest.Assignee.Name,
		milestone,
		strings.Join(mergeRequest.Labels, ", "),
		internal.SweepMarkdownComment(mergeRequest.Description),
	)
	return detial
}
