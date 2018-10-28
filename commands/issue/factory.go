package issue

import (
	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
)

type MethodFactory interface {
	CreateMethod(opt Option, remote *git.RemoteInfo, iid int, factory lab.APIClientFactory) (internal.Method, error)
}

type IssueMethodFactory struct{}

func (c *IssueMethodFactory) CreateMethod(opt Option, remote *git.RemoteInfo, iid int, factory lab.APIClientFactory) (internal.Method, error) {
	if opt.BrowseOption.Browse {
		return &browseMethod{
			opener: &cmd.Browser{},
			remote: remote,
			id:     iid,
		}, nil
	}

	if iid > 0 {
		if opt.CreateUpdateOption.hasEdit() {
			return &updateOnEditorMethod{
				client:   factory.GetIssueClient(),
				opt:      opt.CreateUpdateOption,
				project:  remote.RepositoryFullName(),
				id:       iid,
				editFunc: nil,
			}, nil
		}
		if opt.CreateUpdateOption.hasUpdate() {
			return &updateMethod{
				client:  factory.GetIssueClient(),
				opt:     opt.CreateUpdateOption,
				project: remote.RepositoryFullName(),
				id:      iid,
			}, nil
		}
		return &detailMethod{
			issueClient: factory.GetIssueClient(),
			noteClient:  factory.GetNoteClient(),
			opt:         opt.ShowOption,
			project:     remote.RepositoryFullName(),
			id:          iid,
		}, nil
	}

	// Case of nothing Issue IID
	if opt.CreateUpdateOption.hasEdit() {
		return &createOnEditorMethod{
			issueClient:      factory.GetIssueClient(),
			repositoryClient: factory.GetRepositoryClient(),
			opt:              opt.CreateUpdateOption,
			project:          remote.RepositoryFullName(),
			editFunc:         nil,
		}, nil
	}
	if opt.CreateUpdateOption.hasCreate() {
		return &createMethod{
			client:  factory.GetIssueClient(),
			opt:     opt.CreateUpdateOption,
			project: remote.RepositoryFullName(),
		}, nil
	}
	if opt.ListOption.AllProject {
		return &listAllMethod{
			client: factory.GetIssueClient(),
			opt:    opt.ListOption,
		}, nil
	}

	return &listMethod{
		client:  factory.GetIssueClient(),
		opt:     opt.ListOption,
		project: remote.RepositoryFullName(),
	}, nil
}

type MockMethodFactory struct{}

func (c *MockMethodFactory) CreateMethod(opt Option, remote *git.RemoteInfo, iid int, factory lab.APIClientFactory) (internal.Method, error) {
	return &internal.MockMethod{}, nil
}
