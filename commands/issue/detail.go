package issue

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
	issueClient lab.Issue
	noteClient  lab.Note
	id          int
	project     string
	opt         *ShowOption
}

func (m *detailMethod) Process() (string, error) {
	issue, err := m.issueClient.GetIssue(m.id, m.project)
	if err != nil {
		return "", err
	}
	res := issueDetailOutput(issue)

	if m.opt.NoComment {
		return res, nil
	}

	notes, err := m.noteClient.GetIssueNotes(m.project, m.id, makeListIssueNotesOptions())
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

func makeListIssueNotesOptions() *gitlab.ListIssueNotesOptions {
	return &gitlab.ListIssueNotesOptions{
		Page:    1,
		PerPage: 20,
	}
}

func issueDetailOutput(issue *gitlab.Issue) string {
	base := `%s %s [%s] (created by @%s, %s)
Assignee: %s
Milestone: %s
Labels: %s

%s`

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	var stateColor func(a ...interface{}) string
	if issue.State == "opened" {
		stateColor = color.New(color.FgGreen).SprintFunc()
	} else {
		stateColor = color.New(color.FgRed).SprintFunc()
	}

	milestone := ""
	if issue.Milestone != nil {
		milestone = issue.Milestone.Title
	}

	detial := fmt.Sprintf(base,
		yellow(issue.IID),
		cyan(issue.Title),
		stateColor(issue.State),
		issue.Author.Name,
		issue.CreatedAt.String(),
		issue.Assignee.Name,
		milestone,
		strings.Join(issue.Labels, ", "),
		internal.SweepMarkdownComment(issue.Description),
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
