package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type GitRemote struct {
	Url        string
	Domain     string
	User       string
	Repository string
}

func (r *GitRemote) RepositoryUrl() string {
	params := strings.Join([]string{r.Domain, r.User, r.Repository}, "/")
	return "https://" + params
}

func (r *GitRemote) IssueUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "issues"}, "/")
}

func (r *GitRemote) IssueDetailUrl(issueNo int) string {
	return strings.Join([]string{r.RepositoryUrl(), "issues", fmt.Sprintf("%d", issueNo)}, "/")
}

func (r *GitRemote) BaseUrl() string {
	return "https://" + r.Domain + "/"
}

func NewRemoteUrl(url string) (*GitRemote, error) {
	var (
		otherScheme string
		domain      string
		user        string
		repository  string
	)

	if strings.HasPrefix(url, "ssh") {
		// Case of ssh://git@gitlab.com/lighttiger2505/lab.git
		otherScheme = strings.Split(url, "@")[1]
		otherScheme = strings.TrimSuffix(otherScheme, ".git")

		splitUrl := strings.Split(otherScheme, "/")

		domain = splitUrl[0]
		user = splitUrl[1]
		repository = splitUrl[2]
	} else if strings.HasPrefix(url, "git") {
		// Case of git@gitlab.com/lighttiger2505/lab.git
		otherScheme = strings.Split(url, "@")[1]
		otherScheme = strings.TrimSuffix(otherScheme, ".git")

		splitUrl := strings.Split(otherScheme, ":")
		userRepository := strings.Split(splitUrl[1], "/")

		domain = splitUrl[0]
		user = userRepository[0]
		repository = userRepository[1]
	} else if strings.HasPrefix(url, "https") {
		// Case of https://github.com/lighttiger2505/lab
		otherScheme = strings.Split(url, "//")[1]

		splitUrl := strings.Split(otherScheme, "/")

		domain = splitUrl[0]
		user = splitUrl[1]
		repository = splitUrl[2]
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid remote url: %s", url))
	}

	return &GitRemote{
		Url:        url,
		Domain:     domain,
		User:       user,
		Repository: repository,
	}, nil
}

func gitOutput(name string, args []string) string {
	return gitOutputs(name, args)[0]
}

func gitOutputs(name string, args []string) []string {
	var out = cmdOutput(name, args)
	var outs []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) != "" {
			outs = append(outs, string(line))
		}
	}
	return outs
}

func cmdOutput(name string, args []string) string {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
