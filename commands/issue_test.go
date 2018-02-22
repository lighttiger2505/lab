package commands

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var issues []*gitlabc.Issue = []*gitlabc.Issue{
	&gitlabc.Issue{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo12"},
	&gitlabc.Issue{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo13"},
}

var mockLabIssueClient *gitlab.MockLabClient = &gitlab.MockLabClient{
	MockIssues: func(baseurl, token string, opt *gitlabc.ListIssuesOptions) ([]*gitlabc.Issue, error) {
		return issues, nil
	},
	MockProjectIssues: func(baseurl, token string, opt *gitlabc.ListProjectIssuesOptions, repositoryName string) ([]*gitlabc.Issue, error) {
		return issues, nil
	},
}

func TestIssueCommandRun_Issue(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("token\n")
	c := IssueCommand{
		Ui:           mockUi,
		RemoteFilter: gitlab.NewRemoteFilter(),
		GitClient:    git.NewMockClient(),
		LabClient:    mockLabIssueClient,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUi.ErrorWriter.String())
	}

	got := mockUi.Writer.String()
	want := "#12  Title12\n#13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestIssueCommandRun_ProjectIssue(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("token\n")
	c := IssueCommand{
		Ui:           mockUi,
		RemoteFilter: gitlab.NewRemoteFilter(),
		GitClient:    git.NewMockClient(),
		LabClient:    mockLabIssueClient,
	}

	args := []string{"-a"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUi.ErrorWriter.String())
	}

	got := mockUi.Writer.String()
	want := "#12  namespace/repo12  Title12\n#13  namespace/repo13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestIssueOutput(t *testing.T) {
	got := issueOutput(issues)
	want := []string{
		"#12|namespace/repo12|Title12",
		"#13|namespace/repo13|Title13",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bad return value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestProjectIssueOutput(t *testing.T) {
	got := projectIssueOutput(issues)
	want := []string{
		"#12|Title12",
		"#13|Title13",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bad return value \nwant %#v \ngot  %#v", got, want)
	}
}
