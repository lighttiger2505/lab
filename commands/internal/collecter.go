package internal

import (
	"fmt"
	"strings"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
)

type RemoteCollecter struct {
	UI        ui.Ui
	GitClient git.Client
	Cfg       *config.ConfigV2
}

type GitLabProjectInfo struct {
	Domain  string
	Project string
	Token   string
}

func NewRemoteCollecter(ui ui.Ui, cfg *config.ConfigV2, gitClient git.Client) *RemoteCollecter {
	return &RemoteCollecter{
		UI:        ui,
		Cfg:       cfg,
		GitClient: gitClient,
	}
}

func (c *RemoteCollecter) CollectTarget(project, profile string) (*GitLabProjectInfo, error) {
	pInfo := &GitLabProjectInfo{}
	var err error

	pInfo = c.CollectTargetByDefaultConfig(pInfo)

	pInfo, err = c.CollectTargetByLocalRepository(pInfo)
	if err != nil {
		return nil, err
	}

	pInfo, err = c.CollectTargetByArgs(pInfo, project, profile)
	if err != nil {
		return nil, err
	}

	return pInfo, nil
}

func (c *RemoteCollecter) CollectTargetByDefaultConfig(pInfo *GitLabProjectInfo) *GitLabProjectInfo {
	if c.Cfg.DefalutProfile == "" {
		return pInfo
	}
	profile := c.Cfg.GetDefaultProfile()
	pInfo.Domain = c.Cfg.DefalutProfile
	pInfo.Token = profile.Token

	if profile.DefaultProject == "" {
		return pInfo
	}
	pInfo.Project = profile.DefaultProject

	return pInfo
}

func (c *RemoteCollecter) CollectTargetByLocalRepository(pInfo *GitLabProjectInfo) (*GitLabProjectInfo, error) {
	gitRemotes, err := c.GitClient.RemoteInfos()
	if err != nil {
		return nil, err
	}

	gitlabRemotes := filterHasGitlabDomain(gitRemotes)
	if len(gitlabRemotes) == 0 {
		return nil, fmt.Errorf("Not found gitlab remote repository")
	}
	processedRemotes := excludeDuplicateDomain(gitlabRemotes)

	targetRepo := processedRemotes[0]

	domain := targetRepo.Domain
	if !c.Cfg.HasDomain(domain) {
		c.UI.Message(fmt.Sprintf("Not found this domain [%s].", domain))
		c.Cfg.SetProfile(domain, config.Profile{})
		if err := c.Cfg.Save(); err != nil {
			return nil, err
		}
		c.UI.Message("Saved profile.")
	}

	profile, _ := c.Cfg.GetProfile(domain)
	token := profile.Token
	if token == "" {
		c.UI.Message(fmt.Sprintf("Not found private token in the domain [%s].", domain))
		token, err := c.UI.Ask("Please enter GitLab private token:")
		if err != nil {
			return nil, fmt.Errorf("cannot read private token, %s", err)
		}

		profile.Token = token
		if err := c.Cfg.Save(); err != nil {
			return nil, err
		}
		c.UI.Message("Saved private Token.")
	}

	pInfo.Domain = domain
	pInfo.Token = token
	pInfo.Project = targetRepo.RepositoryFullName()

	return pInfo, nil
}

func (c *RemoteCollecter) CollectTargetByArgs(pInfo *GitLabProjectInfo, project, profile string) (*GitLabProjectInfo, error) {
	if profile != "" {
		p, err := c.Cfg.GetProfile(profile)
		if err != nil {
			return nil, err
		}
		pInfo.Domain = profile
		pInfo.Token = p.Token
	}

	if project != "" {
		pInfo.Project = project
	}

	return pInfo, nil
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
