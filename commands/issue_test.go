package commands

import (
	"reflect"
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var issues = []*gitlabc.Issue{
	&gitlabc.Issue{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo12"},
	&gitlabc.Issue{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo13"},
}

var mockGitlabIssueClient = &gitlab.MockLabClient{
	MockIssues: func(opt *gitlabc.ListIssuesOptions) ([]*gitlabc.Issue, error) {
		return issues, nil
	},
	MockProjectIssues: func(opt *gitlabc.ListProjectIssuesOptions, repositoryName string) ([]*gitlabc.Issue, error) {
		return issues, nil
	},
}

var mockIssueProvider = &gitlab.MockProvider{
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
		return mockGitlabIssueClient, nil
	},
}

func TestIssueCommandRun_Issue(t *testing.T) {
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
	want := "#12  Title12\n#13  Title13\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestIssueCommandRun_ProjectIssue(t *testing.T) {
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
