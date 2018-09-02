package commands

import (
	"reflect"
	"testing"
	"time"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var createdAt, _ = time.Parse("2006-01-02", "2018-02-14")
var updatedAt, _ = time.Parse("2006-01-02", "2018-03-14")

var issue = &gitlab.Issue{
	IID:   12,
	Title: "Title12",
	State: "State12",
	Assignee: struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	}{
		Name: "AssigneeName",
	},
	Author: struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Name      string     `json:"name"`
		State     string     `json:"state"`
		CreatedAt *time.Time `json:"created_at"`
	}{
		Name: "AuthorName",
	},
	CreatedAt:   &createdAt,
	UpdatedAt:   &updatedAt,
	Description: "Description",
}

var issues = []*gitlab.Issue{
	&gitlab.Issue{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo12"},
	&gitlab.Issue{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo13"},
}

var mockGitlabIssueClient = &lab.MockLabIssueClient{
	MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
		return issue, nil
	},
	MockGetAllProjectIssues: func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
		return issues, nil
	},
	MockGetProjectIssues: func(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
		return issues, nil
	},
	MockCreateIssue: func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
		return issue, nil
	},
	MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
		return issue, nil
	},
}

var mockIssueProvider = &lab.MockProvider{
	MockInit: func() error { return nil },
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			Group:      "group",
			Repository: "repository",
		}, nil
	},
	MockGetIssueClient: func(remote *git.RemoteInfo) (lab.Issue, error) {
		return mockGitlabIssueClient, nil
	},
}

func TestIssueCommandRun_ShowIssue(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
	}

	args := []string{"12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := `#12
Title: Title12
Assignee: AssigneeName
Author: AuthorName
State: State12
CreatedAt: 2018-02-14 00:00:00 +0000 UTC
UpdatedAt: 2018-03-14 00:00:00 +0000 UTC

Description
`

	if got != want {
		t.Fatalf("bad output value \nwant %s \ngot  %s", got, want)
	}
}

func TestIssueCommandRun_ListIssue(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "12  Title12\n13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestIssueCommandRun_ListProjectIssue(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
	}

	args := []string{"--all-project"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "namespace/repo12  12  Title12\nnamespace/repo13  13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestIssueCommandRun_CreateIssue(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
	}

	args := []string{"-i", "title", "-m", "message"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "12\n"

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", got, want)
	}
}

func TestIssueCommandRun_CreateIssueOnEditor(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
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

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", got, want)
	}
}

func TestIssueCommandRun_UpdateIssue(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
	}

	args := []string{"-i", "title", "-m", "message", "12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "12\n"

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", got, want)
	}
}

func TestIssueCommandRun_UpdateIssueOnEditor(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := IssueCommand{
		Ui:       mockUI,
		Provider: mockIssueProvider,
		EditFunc: func(program, file string) error {
			return nil
		},
	}

	args := []string{"-e", "-i", "title", "-m", "message", "12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "12\n"

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", got, want)
	}
}

func TestIssueOutput(t *testing.T) {
	got := issueOutput(issues)
	want := []string{
		"namespace/repo12|12|Title12",
		"namespace/repo13|13|Title13",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bad return value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestProjectIssueOutput(t *testing.T) {
	got := projectIssueOutput(issues)
	want := []string{
		"12|Title12",
		"13|Title13",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bad return value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestEditIssueMessage(t *testing.T) {
	got := editIssueMessage("title", "description")
	want := `title

description
`
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}
