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

type Provider struct {
	UI            ui.Ui
	GitClient     git.Client
	ConfigManager *config.ConfigManager
}

func NewProvider(ui ui.Ui, gitClient git.Client, configManager *config.ConfigManager) *Provider {
	return &Provider{
		UI:            ui,
		GitClient:     gitClient,
		ConfigManager: configManager,
	}
}

func (p *Provider) Init() error {
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

func (p *Provider) GetSpecificRemote(namespace, project string) *git.RemoteInfo {
	domain := p.ConfigManager.GetTopDomain()
	return &git.RemoteInfo{
		Domain:     domain,
		NameSpace:  namespace,
		Repository: project,
	}
}

func (p *Provider) GetCurrentRemote() (*git.RemoteInfo, error) {
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

	if len(gitlabRemotes) == 1 {
		return gitlabRemotes[0], nil
	} else if len(gitlabRemotes) < 1 {
		// Current directory is not git repository
		return nil, fmt.Errorf("Not found gitlab remote repository")
	}

	gitlabRemote := registedDomainRemote(gitlabRemotes, p.ConfigManager.Config.PreferredDomains)
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

func (p *Provider) GetClient(remote *git.RemoteInfo) (Client, error) {
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
	return NewLabClient(client), nil
}

func (p *Provider) selectTargetRemote(remoteInfos []*git.RemoteInfo) (*git.RemoteInfo, error) {
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
	sp := strings.Split(webURL, "/")
	return strings.Join([]string{sp[3], sp[4]}, "/")
}
