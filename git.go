package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type GitRemote struct {
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

func (r *GitRemote) MergeRequestUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "merge_requests"}, "/")
}

func (r *GitRemote) MergeRequestDetailUrl(mergeRequestNo int) string {
	return strings.Join([]string{r.RepositoryUrl(), "merge_requests", fmt.Sprintf("%d", mergeRequestNo)}, "/")
}

func (r *GitRemote) BaseUrl() string {
	return "https://" + r.Domain + "/"
}

func (r *GitRemote) ApiUrl() string {
	params := strings.Join([]string{r.Domain, "api", "v4"}, "/")
	return "https://" + params
}

func (r *GitRemote) FullName() string {
	return strings.ToLower(fmt.Sprintf("%s/%s", r.User, r.Repository))
}

func (r *GitRemote) NamespacedPassEncoding() string {
	return fmt.Sprintf("%s%%2F%s", r.User, r.Repository)
}

func NewRemoteUrl(url string) *GitRemote {
	splitUrl := regexp.MustCompile("/|:|@").Split(url, -1)
	return &GitRemote{
		Repository: strings.TrimSuffix(splitUrl[len(splitUrl)-1], ".git"),
		User:       splitUrl[len(splitUrl)-2],
		Domain:     splitUrl[len(splitUrl)-3],
	}
}

func GitRemotes() ([]GitRemote, error) {
	// Get remote repositorys
	remotes := gitOutputs("git", []string{"remote"})
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}
	// Extract domain, namespace, repository name from git remote url
	var gitRemotes []GitRemote
	for _, remote := range remotes {
		url := gitOutput("git", []string{"remote", "get-url", remote})
		gitRemote := NewRemoteUrl(url)
		gitRemotes = append(gitRemotes, *gitRemote)
	}
	return gitRemotes, nil
}
