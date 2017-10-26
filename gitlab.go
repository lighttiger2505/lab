package main

import (
	"errors"
	"fmt"
	"strings"

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
	if err := ReadConfig(); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed read config: %s", err.Error()))
	}

	// Create client
	client := gitlab.NewClient(nil, GetPrivateToken())
	client.SetBaseURL(gitlabRemote.ApiUrl())

	return client, nil
}
