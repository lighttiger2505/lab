package commands

import (
	"testing"

	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/internal/ui"
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

var mockGitlabUserClinet = &lab.MockUserClient{
	MockProjectUsers: func(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error) {
		return testProjectUsers, nil
	},
	MockUsers: func(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error) {
		return testUsers, nil
	},
}

func TestUserCommandRun(t *testing.T) {
	mockClientFactory := &lab.MockAPIClientFactory{
		MockGetUserClient: func() lab.User {
			return mockGitlabUserClinet
		},
	}
	mockUI := ui.NewMockUi()
	c := UserCommand{
		UI:              mockUI,
		RemoteCollecter: &gitutil.MockCollecter{},
		ClientFactory:   mockClientFactory,
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
	mockClientFactory := &lab.MockAPIClientFactory{
		MockGetUserClient: func() lab.User {
			return mockGitlabUserClinet
		},
	}
	mockUI := ui.NewMockUi()
	c := UserCommand{
		UI:              mockUI,
		RemoteCollecter: &gitutil.MockCollecter{},
		ClientFactory:   mockClientFactory,
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
