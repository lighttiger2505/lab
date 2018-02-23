package commands

import (
	"bytes"
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

var mockLabMergeRequestClient *gitlab.MockLabClient = &gitlab.MockLabClient{
	MockMergeRequest: func(baseurl, token string, opt *gitlabc.ListMergeRequestsOptions) ([]*gitlabc.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockProjectMergeRequest: func(baseurl, token string, opt *gitlabc.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlabc.MergeRequest, error) {
		return mergeRequests, nil
	},
}

func TestMergeRequestCommandRun(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("token\n")
	c := MergeRequestCommand{
		Ui:           mockUi,
		RemoteFilter: gitlab.NewRemoteFilter(),
		GitClient:    git.NewMockClient(),
		LabClient:    mockLabMergeRequestClient,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUi.ErrorWriter.String())
	}

	got := mockUi.Writer.String()
	want := "!12  Title12\n!13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestMergeRequestCommandRun_AllProjectOption(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("token\n")
	c := MergeRequestCommand{
		Ui:           mockUi,
		RemoteFilter: gitlab.NewRemoteFilter(),
		GitClient:    git.NewMockClient(),
		LabClient:    mockLabMergeRequestClient,
	}

	args := []string{"-a"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUi.ErrorWriter.String())
	}

	got := mockUi.Writer.String()
	want := "!12  Title12\n!13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}
