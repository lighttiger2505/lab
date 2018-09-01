package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var testPipelines = gitlab.PipelineList{
	struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
		Ref    string `json:"ref"`
		Sha    string `json:"sha"`
	}{
		ID:     1,
		Status: "status1",
		Ref:    "ref1",
		Sha:    "sha1",
	},
	struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
		Ref    string `json:"ref"`
		Sha    string `json:"sha"`
	}{
		ID:     2,
		Status: "status2",
		Ref:    "ref2",
		Sha:    "sha2",
	},
}

var mockGitlabPipelineClient = &lab.MockLabClient{
	MockProjectPipelines: func(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
		return testPipelines, nil
	},
}

var mockPipelineProvider = &lab.MockProvider{
	MockInit: func() error { return nil },
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			Group:      "group",
			Repository: "repository",
		}, nil
	},
	MockGetClient: func(remote *git.RemoteInfo) (lab.Client, error) {
		return mockGitlabPipelineClient, nil
	},
}

func TestPipelineCommandRun(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := PipelineCommand{
		UI:       mockUI,
		Provider: mockPipelineProvider,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "1  status1  ref1  sha1\n2  status2  ref2  sha2\n"

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", want, got)
	}
}
