package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type RemoteInfo struct {
	Domain     string
	NameSpace  string
	Repository string
}

func (r *RemoteInfo) RepositoryUrl() string {
	params := strings.Join([]string{r.Domain, r.NameSpace, r.Repository}, "/")
	return "https://" + params
}

func (r *RemoteInfo) IssueUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "issues"}, "/")
}

func (r *RemoteInfo) IssueDetailUrl(issueNo int) string {
	return strings.Join([]string{r.RepositoryUrl(), "issues", fmt.Sprintf("%d", issueNo)}, "/")
}

func (r *RemoteInfo) MergeRequestUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "merge_requests"}, "/")
}

func (r *RemoteInfo) MergeRequestDetailUrl(mergeRequestNo int) string {
	return strings.Join([]string{r.RepositoryUrl(), "merge_requests", fmt.Sprintf("%d", mergeRequestNo)}, "/")
}

func (r *RemoteInfo) BaseUrl() string {
	return "https://" + r.Domain + "/"
}

func (r *RemoteInfo) ApiUrl() string {
	params := strings.Join([]string{r.Domain, "api", "v4"}, "/")
	return "https://" + params
}

func (r *RemoteInfo) FullName() string {
	return strings.ToLower(fmt.Sprintf("%s/%s", r.NameSpace, r.Repository))
}

func (r *RemoteInfo) NamespacedPassEncoding() string {
	return fmt.Sprintf("%s%%2F%s", r.NameSpace, r.Repository)
}

func NewRemoteInfo(url string) *RemoteInfo {
	splitUrl := regexp.MustCompile("/|:|@").Split(url, -1)
	return &RemoteInfo{
		Repository: strings.TrimSuffix(splitUrl[len(splitUrl)-1], ".git"),
		NameSpace:  splitUrl[len(splitUrl)-2],
		Domain:     splitUrl[len(splitUrl)-3],
	}
}

func GitRemotes() ([]RemoteInfo, error) {
	// Get remote repositorys
	remotes := gitOutputs("git", []string{"remote"})
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}
	// Extract domain, namespace, repository name from git remote url
	var remoteInfos []RemoteInfo
	for _, remote := range remotes {
		url := gitOutput("git", []string{"remote", "get-url", remote})
		remoteInfo := NewRemoteInfo(url)
		remoteInfos = append(remoteInfos, *remoteInfo)
	}
	return remoteInfos, nil
}
