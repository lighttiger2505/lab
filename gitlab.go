package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/xanzy/go-gitlab"
)

func GitlabRemote(ui cli.Ui, config *Config) (*RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := GitRemotes()
	if err != nil {
		return nil, err
	}
	// Filtering only gitlab remote info
	gitlabRemotes := filterGitlab(gitRemotes)

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

func filterGitlab(remoteInfos []RemoteInfo) []RemoteInfo {
	var gitlabRemotes []RemoteInfo
	for _, remoteInfo := range remoteInfos {
		if strings.HasPrefix(remoteInfo.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, remoteInfo)
		}
	}
	return gitlabRemotes
}

func selectUseRemote(ui cli.Ui, gitlabRemotes []RemoteInfo, config *Config) (*RemoteInfo, error) {
	// Search for remote repositorie whose selection is prioritized in the config
	var gitlabRemote *RemoteInfo
	gitlabRemote = hasPriorityRemote(gitlabRemotes, config)
	if gitlabRemote == nil {
		// Get remote repository selected by user input
		var err error
		gitlabRemote, err = inputUseRemote(ui, gitlabRemotes, config)
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

func hasPriorityRemote(remoteInfos []RemoteInfo, config *Config) *RemoteInfo {
	for _, domain := range config.Repositorys {
		for _, remoteInfo := range remoteInfos {
			if domain == remoteInfo.Domain {
				return &remoteInfo
			}
		}
	}
	return nil
}

func inputUseRemote(ui cli.Ui, remoteInfos []RemoteInfo, config *Config) (*RemoteInfo, error) {
	// Receive number of the domain of the remote repository to be searched from stdin
	ui.Info("That repository existing multi gitlab remote repository.")
	for i, remoteInfo := range remoteInfos {
		ui.Info(fmt.Sprintf("%d) %s", i+1, remoteInfo.Domain))
	}
	fmt.Print("Please choice target domain :")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()

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

func GitlabClient(ui cli.Ui, gitlabRemote *RemoteInfo, config *Config) (*gitlab.Client, error) {
	token, err := getPrivateToken(gitlabRemote.Domain, config)
	if err != nil {
		return nil, fmt.Errorf("Failed getting private token. %s", err.Error())
	}

	// Create client
	client := gitlab.NewClient(nil, token)
	apiURL := gitlabRemote.ApiUrl()
	if err := client.SetBaseURL(gitlabRemote.ApiUrl()); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s, %s", apiURL, err.Error())
	}
	return client, nil
}

func getPrivateToken(domain string, config *Config) (string, error) {
	token := ""
	for _, mapItem := range config.Tokens {
		if mapItem.Key.(string) == domain {
			token = mapItem.Value.(string)
		}
	}

	if token == "" {
		fmt.Print("Please input GitLab private token :")
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		token = stdin.Text()

		config.AddToken(domain, token)
		if err := config.Write(); err != nil {
			return "", fmt.Errorf("Failed update config of private token. %s", err.Error())
		}
	}
	return token, nil
}
