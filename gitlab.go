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

func FilterGitlabRemote(ui cli.Ui, remoteInfos []RemoteInfo, config *Config) (*RemoteInfo, error) {
	var gitlabRemotes []RemoteInfo
	for _, remoteInfo := range remoteInfos {
		if strings.HasPrefix(remoteInfo.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, remoteInfo)
		}
	}

	if len(gitlabRemotes) == 1 {
		gitLabRemote := gitlabRemotes[0]
		return &gitLabRemote, nil
	} else if len(gitlabRemotes) > 1 {
		priorityRemote := PriorityRemote(gitlabRemotes, config)
		if priorityRemote != nil {
			return priorityRemote, nil
		} else {
			gitLabRemote, err := ChoiseGitlabRemote(ui, gitlabRemotes, config)
			if err != nil {
				return nil, fmt.Errorf("Failed choise gitlab remote. %v", err.Error())
			}
			ui.Info(fmt.Sprintf("Choised gitlab remote. %s", gitLabRemote.Domain))

			config.AddRepository(gitLabRemote.Domain)
			if err := config.Write(); err != nil {
				return nil, fmt.Errorf("Failed update config of repository priority. %v", err.Error())
			}
			return gitLabRemote, nil
		}
	} else {
		return nil, errors.New("Not a cloned repository from gitlab.")
	}
}

func PriorityRemote(remoteInfos []RemoteInfo, config *Config) *RemoteInfo {
	var priorityRemote RemoteInfo
	for _, repository := range config.Repositorys {
		for _, remoteInfo := range remoteInfos {
			if repository == remoteInfo.Domain {
				priorityRemote = remoteInfo
			}
		}
	}
	return &priorityRemote
}

func ChoiseGitlabRemote(ui cli.Ui, remoteInfos []RemoteInfo, config *Config) (*RemoteInfo, error) {
	fmt.Println("That repository existing multi gitlab remote url.")
	for i, remoteInfo := range remoteInfos {
		ui.Info(fmt.Sprintf("%d) %s", i+1, remoteInfo.Domain))
	}

	ui.Info("Please choice target domain :")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()

	choiceNumber, err := strconv.Atoi(text)
	if err != nil {
		return nil, fmt.Errorf("Failed parse number. %v", err.Error())
	}
	if choiceNumber < 1 {
		return nil, fmt.Errorf("Invalid numver. %d", choiceNumber)
	} else if choiceNumber > len(remoteInfos) {
		return nil, fmt.Errorf("Invalid numver. %d", choiceNumber)
	}
	gitLabRemote := remoteInfos[choiceNumber-1]
	return &gitLabRemote, nil
}

func GitlabRemote(ui cli.Ui, config *Config) (*RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := GitRemotes()
	if err != nil {
		return nil, err
	}
	// Filter gitlab remote url only
	gitlabRemote, err := FilterGitlabRemote(ui, gitRemotes, config)
	if err != nil {
		return nil, err
	}
	return gitlabRemote, nil
}

func GitlabClient(ui cli.Ui, gitlabRemote *RemoteInfo, config *Config) (*gitlab.Client, error) {
	token := ""
	for _, mapItem := range config.Tokens {
		if mapItem.Key.(string) == gitlabRemote.Domain {
			token = mapItem.Value.(string)
		}
	}

	if token == "" {
		fmt.Print("Please input GitLab private token :")
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		token = stdin.Text()

		config.AddToken(gitlabRemote.Domain, token)
		if err := config.Write(); err != nil {
			return nil, fmt.Errorf("Failed update config of private token. %v", err.Error())
		}
	}

	// Create client
	client := gitlab.NewClient(nil, token)
	apiUrl := gitlabRemote.ApiUrl()
	if err := client.SetBaseURL(gitlabRemote.ApiUrl()); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s, %v", apiUrl, err.Error())
	}
	return client, nil
}
