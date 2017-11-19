package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/lab"
	"github.com/xanzy/go-gitlab"
)

func GitlabRemote(ui lab.Ui, config *Config) (*RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := GitRemotes()
	if err != nil {
		return nil, err
	}
	// Filtering only gitlab remote info
	gitlabRemotes := filterHasGitlabDomain(gitRemotes)

	// Filter gitlab remote url only
	var gitlabRemote *RemoteInfo
	if len(gitlabRemotes) == 1 {
		gitlabRemote = &gitlabRemotes[0]
	} else if len(gitlabRemotes) > 1 {
		var err error
		gitlabRemote, err = selectUseRemote(ui, gitlabRemotes, config)
		if err != nil {
			return nil, fmt.Errorf("Failed select multi remote repository. %v", err.Error())
		}
	} else {
		return nil, errors.New("Not a cloned repository from gitlab")
	}
	return gitlabRemote, nil
}

func filterHasGitlabDomain(remoteInfos []RemoteInfo) []RemoteInfo {
	var gitlabRemotes []RemoteInfo
	for _, remoteInfo := range remoteInfos {
		if strings.HasPrefix(remoteInfo.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, remoteInfo)
		}
	}
	return gitlabRemotes
}

func selectUseRemote(ui lab.Ui, gitlabRemotes []RemoteInfo, config *Config) (*RemoteInfo, error) {
	// Search for remote repositorie whose selection is prioritized in the config
	var gitlabRemote *RemoteInfo
	gitlabRemote = hasPriorityRemote(gitlabRemotes, config.PreferredDomains)
	if gitlabRemote == nil {
		// Get remote repository selected by user input
		var err error
		gitlabRemote, err = inputUseRemote(ui, gitlabRemotes)
		if err != nil {
			return nil, fmt.Errorf("Failed choise gitlab remote. %v", err.Error())
		}

		// Add selected remote repository to config
		config.AddRepository(gitlabRemote.Domain)
		if err := config.Write(); err != nil {
			return nil, fmt.Errorf("Failed update config of repository priority. %v", err.Error())
		}
	}
	return gitlabRemote, nil
}

func hasPriorityRemote(remoteInfos []RemoteInfo, preferredDomains []string) *RemoteInfo {
	for _, preferredDomain := range preferredDomains {
		for _, remoteInfo := range remoteInfos {
			if preferredDomain == remoteInfo.Domain {
				return &remoteInfo
			}
		}
	}
	return nil
}

func inputUseRemote(ui lab.Ui, remoteInfos []RemoteInfo) (*RemoteInfo, error) {
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

func GitlabClient(ui lab.Ui, gitlabRemote *RemoteInfo, config *Config) (*gitlab.Client, error) {
	token, err := getPrivateToken(ui, gitlabRemote.Domain, config)
	if err != nil {
		return nil, fmt.Errorf("Failed getting private token. %s", err.Error())
	}

	// Create client
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(gitlabRemote.ApiUrl()); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}
	return client, nil
}

func getPrivateToken(ui lab.Ui, domain string, config *Config) (string, error) {
	token := ""
	for _, mapItem := range config.Tokens {
		if mapItem.Key.(string) == domain {
			token = mapItem.Value.(string)
		}
	}

	if token == "" {
		token, err := ui.Ask("Please input GitLab private token :")
		if err != nil {
			return "", fmt.Errorf("Failed input private token. %s", err.Error())
		}

		config.AddToken(domain, token)
		if err := config.Write(); err != nil {
			return "", fmt.Errorf("Failed update config of private token. %s", err.Error())
		}
	}
	return token, nil
}
