package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func FilterGitlabRemote(gitRemotes []GitRemote) (*GitRemote, error) {
	var gitlabRemotes []GitRemote
	for _, gitRemote := range gitRemotes {
		if strings.HasPrefix(gitRemote.Domain, "gitlab") {
			gitlabRemotes = append(gitlabRemotes, gitRemote)
		}
	}

	var gitLabRemote GitRemote
	if len(gitlabRemotes) > 0 {
		gitLabRemote = gitlabRemotes[0]
	} else {
		return nil, errors.New("Not a cloned repository from gitlab.")
	}
	return &gitLabRemote, nil
}

func GitlabRemote() (*GitRemote, error) {
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

func ProjectId(client *gitlab.Client, gitlabRemote *GitRemote) (int, error) {
	// Search projects
	listProjectOptions := &gitlab.ListProjectsOptions{Search: gitlab.String(gitlabRemote.Repository)}
	projects, _, err := client.Projects.ListProjects(listProjectOptions)
	if err != nil {
		return -1, err
	}

	// Get project id
	projectId := -1
	for _, project := range projects {
		fullName := strings.ToLower(strings.Replace(project.NameWithNamespace, " ", "", -1))
		if fullName == gitlabRemote.FullName() {
			projectId = project.ID
		}
	}
	if projectId == -1 {
		return -1, errors.New(fmt.Sprintf("Failed match Namespace/Project: %s", gitlabRemote.FullName()))
	}
	return projectId, nil
}

func GitlabClient(gitlabRemote *GitRemote) (*gitlab.Client, error) {
	// Read config file
	viper.SetConfigName(".labconfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.lab")
	if err := viper.ReadInConfig(); err != nil {
		if err := CreateConfig(); err != nil {
			return nil, errors.New(fmt.Sprintf("Failed create config file: %s", err.Error()))
		}

		if err := viper.ReadInConfig(); err != nil {
			return nil, errors.New(fmt.Sprintf("Failed read config file: %s", err.Error()))
		}
	}
	privateToken := viper.GetString("private_token")

	// Create client
	client := gitlab.NewClient(nil, privateToken)
	client.SetBaseURL(gitlabRemote.ApiUrl())

	return client, nil
}

func CreateConfig() error {
	dir, err := homedir.Dir()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed get home dir: %s", err.Error()))
	}

	file, err := os.Create(fmt.Sprintf("%s/.labconfig.yml", dir))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed create config file: %s", err.Error()))
	}
	defer file.Close()

	fmt.Print("Plase input GitLab private token :")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	_, err = file.Write([]byte(fmt.Sprintf("private_token: %s", stdin.Text())))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed write config file: %s", err.Error()))
	}

	return nil
}
