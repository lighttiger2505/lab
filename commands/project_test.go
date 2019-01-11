package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var testProjects = []*gitlab.Project{
	&gitlab.Project{
		Name: "name1",
		Namespace: &gitlab.ProjectNamespace{
			Name: "namespace1",
		},
		Description: "description1\ndescription1",
	},
	&gitlab.Project{
		Name: "name2",
		Namespace: &gitlab.ProjectNamespace{
			Name: "namespace2",
		},
		Description: "description2\ndescription2",
	},
}

var mockProjectClient = &lab.MockProjectClient{
	MockProjects: func(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
		return testProjects, nil
	},
}

var mockProjectProvider = &lab.MockProvider{
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			Group:      "group",
			Repository: "repository",
		}, nil
	},
}

func TestProjectCommandRun(t *testing.T) {
	mockClientFactory := &lab.MockAPIClientFactory{
		MockGetProjectClient: func() lab.Project {
			return mockProjectClient
		},
	}
	mockUI := ui.NewMockUi()
	c := ProjectCommand{
		UI:              mockUI,
		RemoteCollecter: &gitutil.MockCollecter{},
		ClientFactory:   mockClientFactory,
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
