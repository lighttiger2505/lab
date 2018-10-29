package pipeline

import (
	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
)

type MethodFactory interface {
	CreateMethod(opt Option, remote *git.RemoteInfo, iid int, factory lab.APIClientFactory) internal.Method
}

type PipelineMethodFacotry struct{}

func (c *PipelineMethodFacotry) CreateMethod(opt Option, remote *git.RemoteInfo, iid int, factory lab.APIClientFactory) internal.Method {
	if opt.BrowseOption.Browse {
		return &browseMethod{
			opener: &cmd.Browser{},
			remote: remote,
			id:     iid,
		}
	}

	if iid > 0 {
		return &listJobMethod{
			client:  factory.GetPipelineClient(),
			opt:     opt.ListOption,
			project: remote.RepositoryFullName(),
			id:      iid,
		}
	}

	return &listMethod{
		client:  factory.GetPipelineClient(),
		opt:     opt.ListOption,
		project: remote.RepositoryFullName(),
	}
}

type MockMethodFactory struct{}

func (c *MockMethodFactory) CreateMethod(opt Option, remote *git.RemoteInfo, iid int, factory lab.APIClientFactory) internal.Method {
	return &internal.MockMethod{}
}
