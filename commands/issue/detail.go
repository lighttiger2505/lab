package issue

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

type detailMethod struct {
	internal.Method
	client  lab.Issue
	id      int
	project string
}

func (m *detailMethod) Process() (string, error) {
	issue, err := m.client.GetIssue(m.id, m.project)
	if err != nil {
		return "", err
	}
	return issueDetailOutput(issue), nil
}

func issueDetailOutput(issue *gitlab.Issue) string {
	base := `#%d
Title: %s
Assignee: %s
Author: %s
State: %s
CreatedAt: %s
UpdatedAt: %s

%s`
	detial := fmt.Sprintf(
		base,
		issue.IID,
		issue.Title,
		issue.Assignee.Name,
		issue.Author.Name,
		issue.State,
		issue.CreatedAt.String(),
		issue.UpdatedAt.String(),
		issue.Description,
	)
	return detial
}
