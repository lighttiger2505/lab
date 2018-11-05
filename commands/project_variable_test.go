package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var mockInit = func() error {
	return nil
}

var mockCurrentRemote = func() (*git.RemoteInfo, error) {
	return &git.RemoteInfo{
		Domain:     "domain",
		Group:      "group",
		Repository: "repository",
	}, nil
}

func TestProjectVariableCommand_Run_List(t *testing.T) {
	// Mocking interfaceis
	sampleProjectVariables := []*gitlab.ProjectVariable{
		&gitlab.ProjectVariable{Key: "foo", Value: "bar"},
		&gitlab.ProjectVariable{Key: "hoge", Value: "soge"},
	}
	mockClient := &lab.MockProjectVariableClient{
		MockGetVariables: func(repositoryName string) ([]*gitlab.ProjectVariable, error) {
			return sampleProjectVariables, nil
		},
	}
	mockProvider := &lab.MockProvider{
		MockGetCurrentRemote: mockCurrentRemote,
		MockGetProjectVariableClient: func(remote *git.RemoteInfo) (lab.ProjectVariable, error) {
			return mockClient, nil
		},
	}
	mockUI := ui.NewMockUi()
	c := ProjectVariableCommand{
		UI:       mockUI,
		Provider: mockProvider,
	}

	// Do command
	args := []string{""}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	// Assertion
	got := mockUI.Writer.String()
	want := "foo   bar\nhoge  soge\n"

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestProjectVariableCommand_Run_Create(t *testing.T) {
	// Mocking interfaceis
	sampleProjectVariable := &gitlab.ProjectVariable{Key: "foo", Value: "bar"}
	mockClient := &lab.MockProjectVariableClient{
		MockCreateVariable: func(repositoryName string, opt *gitlab.CreateVariableOptions) (*gitlab.ProjectVariable, error) {
			return sampleProjectVariable, nil
		},
	}
	mockProvider := &lab.MockProvider{
		MockGetCurrentRemote: mockCurrentRemote,
		MockGetProjectVariableClient: func(remote *git.RemoteInfo) (lab.ProjectVariable, error) {
			return mockClient, nil
		},
	}
	mockUI := ui.NewMockUi()
	c := ProjectVariableCommand{
		UI:       mockUI,
		Provider: mockProvider,
	}

	// Do command
	args := []string{"-a", "foo", "bar"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	// Assertion
	got := mockUI.Writer.String()
	want := ""

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestProjectVariableCommand_Run_Update(t *testing.T) {
	// Mocking interfaceis
	sampleProjectVariable := &gitlab.ProjectVariable{Key: "foo", Value: "bar"}
	mockClient := &lab.MockProjectVariableClient{
		MockUpdateVariable: func(repositoryName string, key string, opt *gitlab.UpdateVariableOptions) (*gitlab.ProjectVariable, error) {
			return sampleProjectVariable, nil
		},
	}
	mockProvider := &lab.MockProvider{
		MockGetCurrentRemote: mockCurrentRemote,
		MockGetProjectVariableClient: func(remote *git.RemoteInfo) (lab.ProjectVariable, error) {
			return mockClient, nil
		},
	}
	mockUI := ui.NewMockUi()
	c := ProjectVariableCommand{
		UI:       mockUI,
		Provider: mockProvider,
	}

	// Do command
	args := []string{"-u", "foo", "bar"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	// Assertion
	got := mockUI.Writer.String()
	want := ""

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestProjectVariableCommand_Run_Remove(t *testing.T) {
	// Mocking interfaceis
	mockClient := &lab.MockProjectVariableClient{
		MockRemoveVariable: func(repositoryName string, key string) error {
			return nil
		},
	}
	mockProvider := &lab.MockProvider{
		MockGetCurrentRemote: mockCurrentRemote,
		MockGetProjectVariableClient: func(remote *git.RemoteInfo) (lab.ProjectVariable, error) {
			return mockClient, nil
		},
	}
	mockUI := ui.NewMockUi()
	c := ProjectVariableCommand{
		UI:       mockUI,
		Provider: mockProvider,
	}

	// Do command
	args := []string{"-d", "foo"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	// Assertion
	got := mockUI.Writer.String()
	want := ""

	if got != want {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}
