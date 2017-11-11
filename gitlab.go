package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func FilterGitlabRemote(remoteInfos []RemoteInfo, config *Config) (*RemoteInfo, error) {
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
			gitLabRemote, err := ChoiseGitlabRemote(gitlabRemotes, config)
			if err != nil {
				return nil, fmt.Errorf("Failed choise gitlab remote. %v", err.Error())
			}
			fmt.Println(fmt.Sprintf("Choised gitlab remote. %s", gitLabRemote.Domain))

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

func ChoiseGitlabRemote(remoteInfos []RemoteInfo, config *Config) (*RemoteInfo, error) {
	fmt.Println("That repository existing multi gitlab remote url.")
	for i, remoteInfo := range remoteInfos {
		fmt.Println(fmt.Sprintf("%d) %s", i+1, remoteInfo.Domain))
	}

	fmt.Print("Please choice target domain :")
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

func GitlabRemote(config *Config) (*RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := GitRemotes()
	if err != nil {
		return nil, err
	}
	// Filter gitlab remote url only
	gitlabRemote, err := FilterGitlabRemote(gitRemotes, config)
	if err != nil {
		return nil, err
	}
	return gitlabRemote, nil
}

func GitlabClient(gitlabRemote *RemoteInfo) (*gitlab.Client, error) {
	c, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed read config: %s", err.Error())
	}

	var token string
	for _, mapItem := range *c.Tokens {
		if mapItem.Key.(string) == gitlabRemote.Domain {
			token = mapItem.Value.(string)
		}
	}

	// Create client
	client := gitlab.NewClient(nil, token)
	client.SetBaseURL(gitlabRemote.ApiUrl())
	return client, nil
}
