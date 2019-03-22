package mr

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/gitutil"
	gitlab "github.com/xanzy/go-gitlab"
)

type createMethod struct {
	internal.Method
	client api.MergeRequest
	opt    *CreateUpdateOption
	pInfo  *gitutil.GitLabProjectInfo
}

func (m *createMethod) Process() (string, error) {
	// Get source branch. current branch from local repository when non specific flags
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return "", err
	}
	if m.opt.SourceBranch != "" {
		// TODO Checking source branch exitst
		currentBranch = m.opt.SourceBranch
	}

	// Do create merge request
	mergeRequest, err := m.client.CreateMergeRequest(
		makeCreateMergeRequestOption(m.opt, m.opt.Title, m.opt.Message, currentBranch, m.pInfo),
		m.pInfo.Project,
	)
	if err != nil {
		return "", err
	}

	// Print created merge request id
	return fmt.Sprintf("%d", mergeRequest.IID), nil
}

type createOnEditorMethod struct {
	internal.Method
	client           api.MergeRequest
	repositoryClient api.Repository
	opt              *CreateUpdateOption
	pInfo            *gitutil.GitLabProjectInfo
	editFunc         func(program, file string) error
}

const templateDir = ".gitlab/merge_request_templates"

func (m *createOnEditorMethod) Process() (string, error) {
	templateFilename := m.opt.Template
	var template string
	if templateFilename != "" {
		filename := templateDir + "/" + templateFilename
		res, err := m.repositoryClient.GetFile(
			m.pInfo.Project,
			filename,
			makeMergeRequestTemplateOption(),
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
		"MERGE_REQUEST",
		internal.EditContents(title, message),
		m.editFunc,
	)
	if err != nil {
		return "", err
	}

	// Get source branch. current branch from local repository when non specific flags
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return "", err
	}
	if m.opt.SourceBranch != "" {
		// TODO Checking source branch exitst
		currentBranch = m.opt.SourceBranch
	}

	// Do create merge request
	mergeRequest, err := m.client.CreateMergeRequest(
		makeCreateMergeRequestOption(m.opt, title, message, currentBranch, m.pInfo),
		m.pInfo.Project,
	)
	if err != nil {
		return "", err
	}

	// Print created merge request id
	return fmt.Sprintf("%d", mergeRequest.IID), nil
}

func makeCreateMergeRequestOption(opt *CreateUpdateOption, title, description, branch string, pInfo *gitutil.GitLabProjectInfo) *gitlab.CreateMergeRequestOptions {
	createMergeRequestOption := &gitlab.CreateMergeRequestOptions{
		Title:           gitlab.String(title),
		Description:     gitlab.String(description),
		SourceBranch:    gitlab.String(branch),
		TargetBranch:    gitlab.String(opt.TargetBranch),
		TargetProjectID: nil,
	}
	assigneeID := opt.getAssigneeID(pInfo.Profile)
	if assigneeID != 0 {
		createMergeRequestOption.AssigneeID = gitlab.Int(assigneeID)
	}
	if opt.MilestoneID != 0 {
		createMergeRequestOption.MilestoneID = gitlab.Int(opt.MilestoneID)
	}
	return createMergeRequestOption
}

func makeMergeRequestTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
	}
	return opt
}
