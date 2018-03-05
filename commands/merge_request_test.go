package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var mergeRequests []*gitlabc.MergeRequest = []*gitlabc.MergeRequest{
	&gitlabc.MergeRequest{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo12"},
	&gitlabc.MergeRequest{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo13"},
}

var mockGitlabMergeRequestClient *gitlab.MockLabClient = &gitlab.MockLabClient{
	MockMergeRequest: func(opt *gitlabc.ListMergeRequestsOptions) ([]*gitlabc.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockProjectMergeRequest: func(opt *gitlabc.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlabc.MergeRequest, error) {
		return mergeRequests, nil
	},
}

var mockMergeRequestProvider = &gitlab.MockProvider{
	MockInit: func() error { return nil },
	MockGetSpecificRemote: func(namespace, project string) *git.RemoteInfo {
		return &git.RemoteInfo{
			Domain:     "domain",
			NameSpace:  "namespace",
			Repository: "repository",
		}
	},
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			NameSpace:  "namespace",
			Repository: "repository",
		}, nil
	},
	MockGetClient: func(remote *git.RemoteInfo) (gitlab.Client, error) {
		return mockGitlabMergeRequestClient, nil
	},
}

func TestMergeRequestCommandRun(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12  Title12\n!13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_AllProjectOption(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"-a"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12  namespace/repo12  Title12\n!13  namespace/repo13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}
