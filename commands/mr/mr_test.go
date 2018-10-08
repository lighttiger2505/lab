package mr

import (
	"testing"
	"time"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var createdAt, _ = time.Parse("2006-01-02", "2018-02-14")
var updatedAt, _ = time.Parse("2006-01-02", "2018-03-14")

var mergeRequest = &gitlab.MergeRequest{
	IID:   12,
	Title: "Title12",
	Assignee: struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Name      string     `json:"name"`
		State     string     `json:"state"`
		CreatedAt *time.Time `json:"created_at"`
	}{
		Name: "AssigneeName",
	},
	Author: struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Name      string     `json:"name"`
		State     string     `json:"state"`
		CreatedAt *time.Time `json:"created_at"`
	}{
		Name: "AuthorName",
	},
	WebURL:      "http://gitlab.jp/namespace/repo12",
	CreatedAt:   &createdAt,
	UpdatedAt:   &updatedAt,
	Description: "Description",
}

var mergeRequests []*gitlab.MergeRequest = []*gitlab.MergeRequest{
	&gitlab.MergeRequest{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo/merge_requests/12"},
	&gitlab.MergeRequest{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo/merge_requests/13"},
}

var mockGitlabMergeRequestClient = &lab.MockLabMergeRequestClient{
	MockGetMergeRequest: func(pid int, repositoryName string) (*gitlab.MergeRequest, error) {
		return mergeRequest, nil
	},
	MockGetAllProjectMergeRequest: func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockGetProjectMargeRequest: func(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockCreateMergeRequest: func(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
		return mergeRequest, nil
	},
	MockUpdateMergeRequest: func(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error) {
		return mergeRequest, nil
	},
}

var mockRepositoryClient = &lab.MockRepositoryClient{
	MockGetFile: func(repositoryName string, filename string, opt *gitlab.GetRawFileOptions) (string, error) {
		return "hogehoge", nil
	},
}

var mockMergeRequestProvider = &lab.MockProvider{
	MockInit: func() error { return nil },
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			Group:      "group",
			Repository: "repository",
		}, nil
	},
	MockGetMergeRequestClient: func(remote *git.RemoteInfo) (lab.MergeRequest, error) {
		return mockGitlabMergeRequestClient, nil
	},
	MockGetRepositoryClient: func(remote *git.RemoteInfo) (lab.Repository, error) {
		return mockRepositoryClient, nil
	},
}

func TestMergeRequestCommandRun_List(t *testing.T) {
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
	want := "12  Title12\n13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_ListAll(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"--all-project"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "namespace/repo  12  Title12\nnamespace/repo  13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_Create(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"-i", "title", "-m", "message"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "12\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_CreateOnEditor(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
		EditFunc: func(program, file string) error {
			return nil
		},
	}

	args := []string{"-e", "-i", "title", "-m", "message"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "12\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_Update(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"-i", "title", "-m", "message", "12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := ""
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_UpdateOnEditor(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
		EditFunc: func(program, file string) error {
			return nil
		},
	}

	args := []string{"-e", "-i", "title", "-m", "message", "12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := ""
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_Show(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := `!12
Title: Title12
Assignee: AssigneeName
Author: AuthorName
CreatedAt: 2018-02-14 00:00:00 +0000 UTC
UpdatedAt: 2018-03-14 00:00:00 +0000 UTC

Description
`
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}
