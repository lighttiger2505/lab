package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var testUsers = []*gitlab.User{
	&gitlab.User{
		ID:       1,
		Username: "username1",
		Name:     "name1",
	},
	&gitlab.User{
		ID:       2,
		Username: "username2",
		Name:     "name2",
	},
}

var testProjectUsers = []*gitlab.ProjectUser{
	&gitlab.ProjectUser{
		ID:       1,
		Username: "username1",
		Name:     "name1",
	},
	&gitlab.ProjectUser{
		ID:       2,
		Username: "username2",
		Name:     "name2",
	},
}

var mockGitlabUserClinet = &lab.MockLabClient{
	MockProjectUsers: func(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error) {
		return testProjectUsers, nil
	},
	MockUsers: func(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error) {
		return testUsers, nil
	},
}

var mockUserProvider = &lab.MockProvider{
	MockInit: func() error { return nil },
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			Group:      "group",
			Repository: "repository",
		}, nil
	},
	MockGetClient: func(remote *git.RemoteInfo) (lab.Client, error) {
		return mockGitlabUserClinet, nil
	},
}

func TestUserCommandRun(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := UserCommand{
		UI:       mockUI,
		Provider: mockUserProvider,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "1  name1  username1\n2  name2  username2\n"

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", want, got)
	}
}

func TestUserCommandRun_AllProject(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := UserCommand{
		UI:       mockUI,
		Provider: mockUserProvider,
	}

	args := []string{"--all-project"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "1  name1  username1\n2  name2  username2\n"

	if got != want {
		t.Fatalf("bad output value \nwant %q \ngot  %q", want, got)
	}
}
