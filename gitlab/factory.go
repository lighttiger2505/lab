package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type APIClientFactory interface {
	Init(url, token string) error
	GetJobClient() Job
	GetIssueClient() Issue
	GetMergeRequestClient() MergeRequest
	GetProjectVariableClient() ProjectVariable
	GetRepositoryClient() Repository
	GetPipelineClient() Pipeline
	GetNoteClient() Note
	GetProjectClient() Project
	GetUserClient() User
	GetLintClient() Lint
	GetRunnerClient() Runner
	GetMilestoneClient() Milestone
	GetBranchClient() Branch
}

type GitlabClientFactory struct {
	gitlabClient *gitlab.Client
}

func NewGitlabClientFactory(url, token string) (APIClientFactory, error) {
	gitlabClient, err := getGitlabClient(url, token)
	if err != nil {
		return nil, err
	}
	factory := &GitlabClientFactory{gitlabClient: gitlabClient}
	return factory, nil
}

func (f *GitlabClientFactory) Init(url, token string) error {
	gitlabClient, err := getGitlabClient(url, token)
	if err != nil {
		return err
	}
	f.gitlabClient = gitlabClient
	return nil
}

func (f *GitlabClientFactory) GetJobClient() Job {
	return NewJobClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetIssueClient() Issue {
	return NewIssueClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetMergeRequestClient() MergeRequest {
	return NewMergeRequestClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetProjectVariableClient() ProjectVariable {
	return NewProjectVariableClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetRepositoryClient() Repository {
	return NewRepositoryClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetNoteClient() Note {
	return NewNoteClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetPipelineClient() Pipeline {
	return NewPipelineClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetProjectClient() Project {
	return NewProjectClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetUserClient() User {
	return NewUserClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetLintClient() Lint {
	return NewLintClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetRunnerClient() Runner {
	return NewRunnerClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetMilestoneClient() Milestone {
	return NewMilestoneClient(f.gitlabClient)
}

func (f *GitlabClientFactory) GetBranchClient() Branch {
	return NewBranchClient(f.gitlabClient)
}

func getGitlabClient(url, token string) (*gitlab.Client, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(url); err != nil {
		return nil, fmt.Errorf("Invalid base url for call GitLab API. %s", err.Error())
	}
	return client, nil
}

type MockAPIClientFactory struct {
	MockGetJobClient             func() Job
	MockGetIssueClient           func() Issue
	MockGetMergeRequestClient    func() MergeRequest
	MockGetProjectVariableClient func() ProjectVariable
	MockGetRepositoryClient      func() Repository
	MockGetNoteClient            func() Note
	MockGetPipelineClient        func() Pipeline
	MockGetProjectClient         func() Project
	MockGetUserClient            func() User
	MockGetLintClient            func() Lint
	MockGetRunnerClient          func() Runner
	MockGetMilestoneClient       func() Milestone
	MockGetBranchClient          func() Branch
}

func (m *MockAPIClientFactory) Init(url, token string) error {
	return nil
}

func (m *MockAPIClientFactory) GetJobClient() Job {
	return m.MockGetJobClient()
}

func (m *MockAPIClientFactory) GetIssueClient() Issue {
	return m.MockGetIssueClient()
}

func (m *MockAPIClientFactory) GetMergeRequestClient() MergeRequest {
	return m.MockGetMergeRequestClient()
}

func (m *MockAPIClientFactory) GetProjectVariableClient() ProjectVariable {
	return m.MockGetProjectVariableClient()
}

func (m *MockAPIClientFactory) GetRepositoryClient() Repository {
	return m.MockGetRepositoryClient()
}

func (m *MockAPIClientFactory) GetPipelineClient() Pipeline {
	return m.MockGetPipelineClient()
}

func (m *MockAPIClientFactory) GetNoteClient() Note {
	return m.MockGetNoteClient()
}

func (m *MockAPIClientFactory) GetProjectClient() Project {
	return m.MockGetProjectClient()
}

func (m *MockAPIClientFactory) GetUserClient() User {
	return m.MockGetUserClient()
}

func (m *MockAPIClientFactory) GetLintClient() Lint {
	return m.MockGetLintClient()
}

func (m *MockAPIClientFactory) GetRunnerClient() Runner {
	return m.MockGetRunnerClient()
}

func (m *MockAPIClientFactory) GetMilestoneClient() Milestone {
	return m.MockGetMilestoneClient()
}

func (m *MockAPIClientFactory) GetBranchClient() Branch {
	return m.MockGetBranchClient()
}
