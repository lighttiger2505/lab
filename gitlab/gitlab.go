package gitlab

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

type Provider interface {
	Init() error
	GetCurrentRemote() (*git.RemoteInfo, error)
	GetAPIToken(remote *git.RemoteInfo) (string, error)
	GetJobClient(remote *git.RemoteInfo) (Job, error)
	GetProjectVariableClient(remote *git.RemoteInfo) (ProjectVariable, error)
	GetRepositoryClient(remote *git.RemoteInfo) (Repository, error)
	GetNoteClient(remote *git.RemoteInfo) (Note, error)
}

type GitlabProvider struct {
	UI            ui.Ui
	GitClient     git.Client
	ConfigManager *config.ConfigManager
}

func NewProvider(ui ui.Ui, gitClient git.Client, configManager *config.ConfigManager) *GitlabProvider {
	return &GitlabProvider{
		UI:            ui,
		GitClient:     gitClient,
		ConfigManager: configManager,
	}
}

func (p *GitlabProvider) Init() error {
	// Load config
	if err := p.ConfigManager.Init(); err != nil {
		return err
	}
	_, err := p.ConfigManager.Load()
	if err != nil {
		return err
	}
	return nil
}

func (p *GitlabProvider) GetCurrentRemote() (*git.RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := p.GitClient.RemoteInfos()
	if err != nil {
		return nil, err
	}

	// Filtering only gitlab remote info
	gitlabRemotes := filterHasGitlabDomain(gitRemotes)
	if err != nil {
		return nil, fmt.Errorf("Failed getting remote info. Error: %v", err.Error())
	}

	if len(gitlabRemotes) == 0 {
		// Current directory is not git repository
		return nil, fmt.Errorf("Not found gitlab remote repository")
	}

	processedRemotes := excludeDuplicateDomain(gitlabRemotes)

	if len(gitlabRemotes) == 1 {
		return gitlabRemotes[0], nil
	}

	gitlabRemote := registedDomainRemote(processedRemotes, p.ConfigManager.Config.PreferredDomains)
	if gitlabRemote == nil {
		// Get remote repository selected by user input
		var err error
		gitlabRemote, err = p.selectTargetRemote(gitlabRemotes)
		if err != nil {
			return nil, fmt.Errorf("Failed choise gitlab remote. %v", err.Error())
		}

		// Add selected remote repository to config
		if err := p.ConfigManager.SavePreferredDomain(gitlabRemote.Domain); err != nil {
			return nil, fmt.Errorf("Failed save preferred domain to config. Error: %v", err.Error())
		}
	}
	return gitlabRemote, nil
}

func excludeDuplicateDomain(remotes []*git.RemoteInfo) []*git.RemoteInfo {
	domainRemotesMap := map[string][]*git.RemoteInfo{}
	for _, remote := range remotes {
		domain := remote.Domain
		domainRemotesMap[domain] = append(domainRemotesMap[domain], remote)
	}

	processedRemotes := []*git.RemoteInfo{}
	for _, v := range domainRemotesMap {
		var tmpRemote = v[0]
		for _, remote := range v {
			if remote.Remote == "origin" {
				tmpRemote = remote
				break
			}
		}
		processedRemotes = append(processedRemotes, tmpRemote)
	}
	return processedRemotes
}

func (p *GitlabProvider) makeGitLabClient(remote *git.RemoteInfo) (*gitlab.Client, error) {
	token := p.ConfigManager.GetTokenOnly(remote.Domain)
	if token == "" {
		token, err := p.UI.Ask("Please input GitLab private token :")
		if err != nil {
			return nil, fmt.Errorf("Failed input private token. %s", err.Error())
		}

		if err := p.ConfigManager.SaveToken(remote.Domain, token); err != nil {
			return nil, fmt.Errorf("Failed update config of private token. %s", err.Error())
		}
	}

	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(remote.ApiUrl()); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}
	return client, nil
}

func (p *GitlabProvider) GetAPIToken(remote *git.RemoteInfo) (string, error) {
	token := p.ConfigManager.GetTokenOnly(remote.Domain)
	if token == "" {
		token, err := p.UI.Ask("Please input GitLab private token :")
		if err != nil {
			return "", fmt.Errorf("Failed input private token. %s", err.Error())
		}

		if err := p.ConfigManager.SaveToken(remote.Domain, token); err != nil {
			return "", fmt.Errorf("Failed update config of private token. %s", err.Error())
		}
	}
	return token, nil
}

func (p *GitlabProvider) GetJobClient(remote *git.RemoteInfo) (Job, error) {
	gitlabClient, err := p.makeGitLabClient(remote)
	if err != nil {
		return nil, err
	}
	return NewJobClient(gitlabClient), nil
}

func (p *GitlabProvider) GetProjectVariableClient(remote *git.RemoteInfo) (ProjectVariable, error) {
	gitlabClient, err := p.makeGitLabClient(remote)
	if err != nil {
		return nil, err
	}
	return NewProjectVariableClient(gitlabClient), nil
}

func (p *GitlabProvider) GetRepositoryClient(remote *git.RemoteInfo) (Repository, error) {
	gitlabClient, err := p.makeGitLabClient(remote)
	if err != nil {
		return nil, err
	}
	return NewRepositoryClient(gitlabClient), nil
}

func (p *GitlabProvider) GetNoteClient(remote *git.RemoteInfo) (Note, error) {
	gitlabClient, err := p.makeGitLabClient(remote)
	if err != nil {
		return nil, err
	}
	return NewNoteClient(gitlabClient), nil
}

func (p *GitlabProvider) selectTargetRemote(remoteInfos []*git.RemoteInfo) (*git.RemoteInfo, error) {
	// Receive number of the domain of the remote repository to be searched from stdin
	p.UI.Message("That repository existing multi gitlab remote repository.")
	for i, remoteInfo := range remoteInfos {
		p.UI.Message(fmt.Sprintf("%d) %s", i+1, remoteInfo.Domain))
	}
	text, err := p.UI.Ask("Please choice target domain :")
	if err != nil {
		return nil, fmt.Errorf("Failed target domain input. %v", err.Error())
	}

	// Check valid number
	choiceNumber, err := strconv.Atoi(text)
	if err != nil {
		return nil, fmt.Errorf("Failed parse number. Error: %s", err.Error())
	}
	if choiceNumber < 1 || choiceNumber > len(remoteInfos) {
		return nil, fmt.Errorf("Invalid number. Input: %d", choiceNumber)
	}

	return remoteInfos[choiceNumber-1], nil
}

func filterHasGitlabDomain(remoteInfos []*git.RemoteInfo) []*git.RemoteInfo {
	var gitlabRemotes []*git.RemoteInfo
	for _, remoteInfo := range remoteInfos {
		if strings.HasPrefix(remoteInfo.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, remoteInfo)
		}
	}
	return gitlabRemotes
}

func registedDomainRemote(remoteInfos []*git.RemoteInfo, resistedDomains []string) *git.RemoteInfo {
	for _, preferredDomain := range resistedDomains {
		for _, remoteInfo := range remoteInfos {
			if preferredDomain == remoteInfo.Domain {
				return remoteInfo
			}
		}
	}
	return nil
}

func ParceRepositoryFullName(webURL string) string {
	splitURL := strings.Split(webURL, "/")[3:]

	subPageWords := []string{
		"issues",
		"merge_requests",
	}
	var subPageIndex int
	for i, word := range splitURL {
		for _, subPageWord := range subPageWords {
			if word == subPageWord {
				subPageIndex = i
			}
		}
	}

	return strings.Join(splitURL[:subPageIndex], "/")
}

type MockProvider struct {
	Provider
	MockInit                     func() error
	MockGetSpecificRemote        func(namespace, project string) *git.RemoteInfo
	MockGetCurrentRemote         func() (*git.RemoteInfo, error)
	MockGetProjectVariableClient func(remote *git.RemoteInfo) (ProjectVariable, error)
	MockGetRepositoryClient      func(remote *git.RemoteInfo) (Repository, error)

	MockGetNoteClient func(remote *git.RemoteInfo) (Note, error)
}

func (m *MockProvider) Init() error {
	return m.MockInit()
}

func (m *MockProvider) GetCurrentRemote() (*git.RemoteInfo, error) {
	return m.MockGetCurrentRemote()
}

func (m *MockProvider) GetAPIToken(remote *git.RemoteInfo) (string, error) {
	return "", nil
}

func (m *MockProvider) GetProjectVariableClient(remote *git.RemoteInfo) (ProjectVariable, error) {
	return m.MockGetProjectVariableClient(remote)
}

func (m *MockProvider) GetRepositoryClient(remote *git.RemoteInfo) (Repository, error) {
	return m.MockGetRepositoryClient(remote)
}

func (m *MockProvider) GetNoteClient(remote *git.RemoteInfo) (Note, error) {
	return m.MockGetNoteClient(remote)
}

func getGitlabClient(url, token string) (*gitlab.Client, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(url); err != nil {
		return nil, fmt.Errorf("Invalid base url for call GitLab API. %s", err.Error())
	}
	return client, nil
}

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
