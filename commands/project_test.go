package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var testProjects = []*gitlabc.Project{
	&gitlabc.Project{
		Name: "name1",
		Namespace: &gitlabc.ProjectNamespace{
			Name: "namespace1",
		},
		Description: "description1\ndescription1",
	},
	&gitlabc.Project{
		Name: "name2",
		Namespace: &gitlabc.ProjectNamespace{
			Name: "namespace2",
		},
		Description: "description2\ndescription2",
	},
}

var mockGitlabProjectClient = &gitlab.MockLabClient{
	MockProjects: func(opt *gitlabc.ListProjectsOptions) ([]*gitlabc.Project, error) {
		return testProjects, nil
	},
}

var mockProjectProvider = &gitlab.MockProvider{
	MockInit: func() error { return nil },
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			NameSpace:  "namespace",
			Repository: "repository",
		}, nil
	},
	MockGetClient: func(remote *git.RemoteInfo) (gitlab.Client, error) {
		return mockGitlabProjectClient, nil
	},
}

func TestProjectCommandRun(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := ProjectCommand{
		UI:       mockUI,
		Provider: mockProjectProvider,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "namespace1/name1  description1description1\nnamespace2/name2  description2description2\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", got, want)
	}
}

func TestRemoveLineBreak(t *testing.T) {
	got := removeLineBreak("123\r\n456\r789\n")
	want := "123456789"
	if got != want {
		t.Fatalf("want %q, but %q:", want, got)
	}
}
