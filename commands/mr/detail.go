package mr

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/internal/api"
	gitlab "github.com/xanzy/go-gitlab"
)

type detailMethod struct {
	internal.Method
	mrClient   api.MergeRequest
	noteClient api.Note
	opt        *ShowOption
	project    string
	id         int
}

func (m *detailMethod) Process() (string, error) {
	// Do get merge request
	mergeRequest, err := m.mrClient.GetMergeRequest(m.id, m.project)
	if err != nil {
		return "", err
	}
	res := outMergeRequestDetail(mergeRequest)

	if m.opt.NoComment {
		return res, nil
	}

	notes, err := m.noteClient.GetMergeRequestNotes(m.project, m.id, makeListMergeRequestNotesOptions())
	if err != nil {
		return "", err
	}
	noteOutputs := make([]string, len(notes))
	for i, note := range notes {
		noteOutputs[i] = noteOutput(note)
	}
	res = res + strings.Join(noteOutputs, "\n")

	return res, nil
}

func makeListMergeRequestNotesOptions() *gitlab.ListMergeRequestNotesOptions {
	listOption := gitlab.ListOptions{
		Page:    1,
		PerPage: 20,
	}
	return &gitlab.ListMergeRequestNotesOptions{
		ListOptions: listOption,
	}
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

func noteOutput(note *gitlab.Note) string {
	base := `
%s (created by @%s, %s)

%s`

	yellow := color.New(color.FgYellow).SprintFunc()
	return fmt.Sprintf(base,
		yellow(fmt.Sprintf("comment %d", note.ID)),
		note.Author.Name,
		note.CreatedAt.String(),
		internal.SweepMarkdownComment(note.Body),
	)
}
