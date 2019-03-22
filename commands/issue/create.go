package issue

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/gitutil"
	gitlab "github.com/xanzy/go-gitlab"
)

func makeCreateIssueOptions(opt *CreateUpdateOption, title, description string, pInfo *gitutil.GitLabProjectInfo) *gitlab.CreateIssueOptions {
	createIssueOption := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(description),
	}
	assigneeID := opt.getAssigneeID(pInfo.Profile)
	if assigneeID != 0 {
		createIssueOption.AssigneeIDs = []int{assigneeID}
	}
	if opt.MilestoneID != 0 {
		createIssueOption.MilestoneID = gitlab.Int(opt.MilestoneID)
	}
	return createIssueOption
}

func makeIssueTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
	}
	return opt
}

type createMethod struct {
	client  api.Issue
	opt     *CreateUpdateOption
	project string
	pInfo   *gitutil.GitLabProjectInfo
}

func (m *createMethod) Process() (string, error) {
	issue, err := m.client.CreateIssue(
		makeCreateIssueOptions(m.opt, m.opt.Title, m.opt.Message, m.pInfo),
		m.project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", issue.IID), nil
}

type createOnEditorMethod struct {
	issueClient      api.Issue
	repositoryClient api.Repository
	opt              *CreateUpdateOption
	editFunc         func(program, file string) error
	pInfo            *gitutil.GitLabProjectInfo
}

const templateDir = ".gitlab/issue_templates"

func (m *createOnEditorMethod) Process() (string, error) {
	templateFilename := m.opt.Template
	var template string
	if templateFilename != "" {
		filename := templateDir + "/" + templateFilename
		res, err := m.repositoryClient.GetFile(
			m.pInfo.Project,
			filename,
			makeIssueTemplateOption(),
		)
		if err != nil {
			return "", err
		}
		template = res
	}

	var title, message string
	title = m.opt.Title
	message = template
	if m.opt.Message != "" {
		message = m.opt.Message
	}

	title, message, err := internal.EditTitleAndDesc(
		"ISSUE",
		internal.EditContents(title, message),
		m.editFunc,
	)
	if err != nil {
		return "", err
	}
	issue, err := m.issueClient.CreateIssue(
		makeCreateIssueOptions(m.opt, title, message, m.pInfo),
		m.pInfo.Project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", issue.IID), nil
}
