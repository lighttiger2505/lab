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

func FilterGitlabRemote(remoteInfos []RemoteInfo) (*RemoteInfo, error) {
	var gitlabRemotes []RemoteInfo
	for _, remoteInfo := range remoteInfos {
		if strings.HasPrefix(remoteInfo.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, remoteInfo)
		}
	}

	var gitLabRemote RemoteInfo
	if len(gitlabRemotes) == 1 {
		gitLabRemote = gitlabRemotes[0]
	} else if len(gitlabRemotes) > 1 {
		fmt.Println("That repository existing multi gitlab remote url.")
		for i, remoteInfo := range gitlabRemotes {
			fmt.Println(fmt.Sprintf("%d) %s", i+1, remoteInfo.Domain))
		}

		fmt.Print(" Please choice target domain :")
		stdin := bufio.NewScanner(os.Stdin)
		stdin.Scan()
		text := stdin.Text()

		choiceNumber, err := strconv.Atoi(text)
		if err != nil {
			return nil, fmt.Errorf("Failed parse number. %v", err.Error())
		}
		if choiceNumber < 1 {
			return nil, fmt.Errorf("Invalid numver. %d", choiceNumber)
		} else if choiceNumber > len(gitlabRemotes) {
			return nil, fmt.Errorf("Invalid numver. %d", choiceNumber)
		}

		gitLabRemote = gitlabRemotes[choiceNumber-1]
	} else {
		return nil, errors.New("Not a cloned repository from gitlab.")
	}
	return &gitLabRemote, nil
}

func GitlabRemote() (*RemoteInfo, error) {
	// Get remote urls
	gitRemotes, err := GitRemotes()
	if err != nil {
		return nil, err
	}
	// Filter gitlab remote url only
	gitlabRemote, err := FilterGitlabRemote(gitRemotes)
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
