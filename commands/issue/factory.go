package issue

import (
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/browse"
	"github.com/lighttiger2505/lab/internal/gitutil"
)

type MethodFactory interface {
	CreateMethod(opt Option, pInfo *gitutil.GitLabProjectInfo, iid int, factory lab.APIClientFactory) internal.Method
}

type IssueMethodFactory struct{}

func (c *IssueMethodFactory) CreateMethod(opt Option, pInfo *gitutil.GitLabProjectInfo, iid int, factory lab.APIClientFactory) internal.Method {
	if opt.BrowseOption.Browse {
		return &browseMethod{
			opener: &browse.Browser{},
			url:    pInfo.SubpageUrl("issues"),
			id:     iid,
		}
	}

	if iid > 0 {
		if opt.CreateUpdateOption.hasEdit() {
			return &updateOnEditorMethod{
				client:   factory.GetIssueClient(),
				opt:      opt.CreateUpdateOption,
				project:  pInfo.Project,
				id:       iid,
				editFunc: nil,
			}
		}
		if opt.CreateUpdateOption.hasUpdate() {
			return &updateMethod{
				client:  factory.GetIssueClient(),
				opt:     opt.CreateUpdateOption,
				project: pInfo.Project,
				id:      iid,
			}
		}
		return &detailMethod{
			issueClient: factory.GetIssueClient(),
			noteClient:  factory.GetNoteClient(),
			opt:         opt.ShowOption,
			project:     pInfo.Project,
			id:          iid,
		}
	}

	// Case of nothing Issue id
	if opt.CreateUpdateOption.hasEdit() {
		return &createOnEditorMethod{
			issueClient:      factory.GetIssueClient(),
			repositoryClient: factory.GetRepositoryClient(),
			opt:              opt.CreateUpdateOption,
			project:          pInfo.Project,
			editFunc:         nil,
		}
	}
	if opt.CreateUpdateOption.hasCreate() {
		return &createMethod{
			client:  factory.GetIssueClient(),
			opt:     opt.CreateUpdateOption,
			project: pInfo.Project,
		}
	}
	if opt.ListOption.AllProject {
		return &listAllMethod{
			client: factory.GetIssueClient(),
			opt:    opt.ListOption,
		}
	}

	return &listMethod{
		client:  factory.GetIssueClient(),
		opt:     opt.ListOption,
		project: pInfo.Project,
	}
}

type MockMethodFactory struct{}

func (c *MockMethodFactory) CreateMethod(opt Option, pInfo *gitutil.GitLabProjectInfo, iid int, factory lab.APIClientFactory) internal.Method {
	return &internal.MockMethod{}
}
