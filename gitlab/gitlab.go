package gitlab

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
	"github.com/xanzy/go-gitlab"
)

type RemoteFilter interface {
	Filter(ui.Ui, *config.Config) (*git.RemoteInfo, error)
}

type GitlabRemoteFilter struct {
}

func (g *GitlabRemoteFilter) Filter(ui ui.Ui, conf *config.Config) (*git.RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := git.GitRemotes()
	if err != nil {
		return nil, err
	}

	// Filtering only gitlab remote info
	gitlabRemotes := filterHasGitlabDomain(gitRemotes)

	// Filter gitlab remote url only
	var gitlabRemote *git.RemoteInfo
	if len(gitlabRemotes) == 1 {
		gitlabRemote = &gitlabRemotes[0]
	} else if len(gitlabRemotes) > 1 {
		var err error
		gitlabRemote, err = selectUseRemote(ui, gitlabRemotes, conf)
		if err != nil {
			return nil, fmt.Errorf("Failed select multi remote repository. %v", err.Error())
		}
	} else {
		// Current directory is not git repository
		return nil, nil
	}
	return gitlabRemote, nil
}

func GitlabRemote(ui ui.Ui, conf *config.Config) (*git.RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := git.GitRemotes()
	if err != nil {
		return nil, err
	}
	// Filtering only gitlab remote info
	gitlabRemotes := filterHasGitlabDomain(gitRemotes)

	// Filter gitlab remote url only
	var gitlabRemote *git.RemoteInfo
	if len(gitlabRemotes) == 1 {
		gitlabRemote = &gitlabRemotes[0]
	} else if len(gitlabRemotes) > 1 {
		var err error
		gitlabRemote, err = selectUseRemote(ui, gitlabRemotes, conf)
		if err != nil {
			return nil, fmt.Errorf("Failed select multi remote repository. %v", err.Error())
		}
	} else {
		// Current directory is not git repository
		return nil, nil
	}
	return gitlabRemote, nil
}

func filterHasGitlabDomain(remoteInfos []git.RemoteInfo) []git.RemoteInfo {
	var gitlabRemotes []git.RemoteInfo
	for _, remoteInfo := range remoteInfos {
		if strings.HasPrefix(remoteInfo.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, remoteInfo)
		}
	}
	return gitlabRemotes
}

func selectUseRemote(ui ui.Ui, gitlabRemotes []git.RemoteInfo, conf *config.Config) (*git.RemoteInfo, error) {
	// Search for remote repositorie whose selection is prioritized in the config
	var gitlabRemote *git.RemoteInfo
	gitlabRemote = hasPriorityRemote(gitlabRemotes, conf.PreferredDomains)
	if gitlabRemote == nil {
		// Get remote repository selected by user input
		var err error
		gitlabRemote, err = inputUseRemote(ui, gitlabRemotes)
		if err != nil {
			return nil, fmt.Errorf("Failed choise gitlab remote. %v", err.Error())
		}

		// Add selected remote repository to config
		conf.AddRepository(gitlabRemote.Domain)
		if err := conf.Write(); err != nil {
			return nil, fmt.Errorf("Failed update config of repository priority. %v", err.Error())
		}
	}
	return gitlabRemote, nil
}

func hasPriorityRemote(remoteInfos []git.RemoteInfo, preferredDomains []string) *git.RemoteInfo {
	for _, preferredDomain := range preferredDomains {
		for _, remoteInfo := range remoteInfos {
			if preferredDomain == remoteInfo.Domain {
				return &remoteInfo
			}
		}
	}
	return nil
}

func inputUseRemote(ui ui.Ui, remoteInfos []git.RemoteInfo) (*git.RemoteInfo, error) {
	// Receive number of the domain of the remote repository to be searched from stdin
	ui.Message("That repository existing multi gitlab remote repository.")
	for i, remoteInfo := range remoteInfos {
		ui.Message(fmt.Sprintf("%d) %s", i+1, remoteInfo.Domain))
	}
	text, err := ui.Ask("Please choice target domain :")
	if err != nil {
		return nil, fmt.Errorf("Failed target domain input. %v", err.Error())
	}

	// Check valid number
	choiceNumber, err := strconv.Atoi(text)
	if err != nil {
		return nil, fmt.Errorf("Failed parse number. %v", err.Error())
	}
	if choiceNumber < 1 || choiceNumber > len(remoteInfos) {
		return nil, fmt.Errorf("Invalid numver. %d", choiceNumber)
	}

	gitLabRemote := &remoteInfos[choiceNumber-1]
	return gitLabRemote, nil
}

func NewGitlabClient(ui ui.Ui, gitlabRemote *git.RemoteInfo, token string) (*gitlab.Client, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(gitlabRemote.ApiUrl()); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}
	return client, nil
}

func ParceRepositoryFullName(webURL string) string {
	sp := strings.Split(webURL, "/")
	return strings.Join([]string{sp[3], sp[4]}, "/")
}

type Client interface {
	Issues(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	ProjectIssues(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MergeRequest(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	ProjectMergeRequest(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
}

type LabClient struct {
	Client
}

func NewLabClient() *LabClient {
	return &LabClient{}
}

func (g *LabClient) Issues(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	issues, _, err := client.Issues.ListIssues(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list issue. %s", err.Error())
	}
	return issues, nil
}

func (g *LabClient) ProjectIssues(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	issues, _, err := client.Issues.ListProjectIssues(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}
	return issues, nil
}

func (g *LabClient) MergeRequest(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	mergeRequests, _, err := client.MergeRequests.ListMergeRequests(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (g *LabClient) ProjectMergeRequest(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

type MockLabClient struct {
	MockIssues              func(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	MockProjectIssues       func(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MockMergeRequest        func(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	MockProjectMergeRequest func(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
}

func (m *MockLabClient) Issues(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	return m.MockIssues(baseurl, token, opt)
}

func (m *MockLabClient) ProjectIssues(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	return m.MockProjectIssues(baseurl, token, opt, repositoryName)
}

func (m *MockLabClient) MergeRequest(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	return m.MockMergeRequest(baseurl, token, opt)
}

func (m *MockLabClient) ProjectMergeRequest(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	return m.MockProjectMergeRequest(baseurl, token, opt, repositoryName)
}
