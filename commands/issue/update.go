package issue

import (
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/internal/api"
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
	if opt.MilestoneID != 0 {
		updateIssueOption.MilestoneID = gitlab.Int(opt.MilestoneID)
	}
	return updateIssueOption
}

type updateMethod struct {
	client  api.Issue
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
	updatedTitle, updatedMessage := getUpdatedTitleAndMessage(issue, m.opt.Title, m.opt.Message)

	// Do update issue
	_, err = m.client.UpdateIssue(
		makeUpdateIssueOption(m.opt, updatedTitle, updatedMessage),
		m.id,
		m.project,
	)
	if err != nil {
		return "", err
	}

	// Return empty value
	return "", nil
}

type updateOnEditorMethod struct {
	internal.Method
	client   api.Issue
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
	updatedTitle, updatedMessage := getUpdatedTitleAndMessage(issue, m.opt.Title, m.opt.Message)

	// Starting editor for edit title and description
	content := editIssueMessage(updatedTitle, updatedMessage)
	title, message, err := editIssueTitleAndDesc(content, m.editFunc)
	if err != nil {
		return "", err
	}

	// Do update issue
	_, err = m.client.UpdateIssue(
		makeUpdateIssueOption(m.opt, title, message),
		m.id,
		m.project,
	)
	if err != nil {
		return "", err
	}

	// Return empty value
	return "", nil
}

func getUpdatedTitleAndMessage(issue *gitlab.Issue, title, message string) (string, string) {
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if title != "" {
		updatedTitle = title
	}
	if message != "" {
		updatedMessage = message
	}
	return updatedTitle, updatedMessage
}
