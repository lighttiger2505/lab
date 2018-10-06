package issue

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

func makeUpdateIssueOption(opt *CreateUpdateOption, title, description string) *gitlab.UpdateIssueOptions {
	updateIssueOption := &gitlab.UpdateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(description),
	}
	if opt.StateEvent != "" {
		updateIssueOption.StateEvent = gitlab.String(opt.StateEvent)
	}
	if opt.AssigneeID != 0 {
		updateIssueOption.AssigneeIDs = []int{opt.AssigneeID}
	}
	return updateIssueOption
}

type updateMethod struct {
	internal.Method
	client  lab.Issue
	opt     *CreateUpdateOption
	project string
	id      int
}

func (m *updateMethod) Process() (string, error) {
	// Getting exist issue
	issue, err := m.client.GetIssue(m.id, m.project)
	if err != nil {
		return "", err
	}

	// Create new title or description
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if m.opt.Title != "" {
		updatedTitle = m.opt.Title
	}
	if m.opt.Message != "" {
		updatedMessage = m.opt.Message
	}

	// Do update issue
	updatedIssue, err := m.client.UpdateIssue(
		makeUpdateIssueOption(m.opt, updatedTitle, updatedMessage),
		m.id,
		m.project,
	)
	if err != nil {
		return "", err
	}

	// Print update Issue IID
	return fmt.Sprintf("%d", updatedIssue.IID), nil
}

type updateOnEditorMethod struct {
	internal.Method
	client   lab.Issue
	opt      *CreateUpdateOption
	project  string
	id       int
	editFunc func(program, file string) error
}

func (m *updateOnEditorMethod) Process() (string, error) {
	// Getting exist issue
	issue, err := m.client.GetIssue(m.id, m.project)
	if err != nil {
		return "", err
	}

	// Create new title or description
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if m.opt.Title != "" {
		updatedTitle = m.opt.Title
	}
	if m.opt.Message != "" {
		updatedMessage = m.opt.Message
	}

	// Starting editor for edit title and description
	content := editIssueMessage(updatedTitle, updatedMessage)
	title, message, err := editIssueTitleAndDesc(content, m.editFunc)
	if err != nil {
		return "", err
	}

	// Do update issue
	updatedIssue, err := m.client.UpdateIssue(
		makeUpdateIssueOption(m.opt, title, message),
		m.id,
		m.project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", updatedIssue.IID), nil
}
