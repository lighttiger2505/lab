package git

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lighttiger2505/lab/utils"
)

type RemoteInfo struct {
	Domain     string
	NameSpace  string
	Repository string
}

func NewRemoteInfo(url string) *RemoteInfo {
	splitUrl := regexp.MustCompile("/|:|@").Split(url, -1)
	return &RemoteInfo{
		Repository: strings.TrimSuffix(splitUrl[len(splitUrl)-1], ".git"),
		NameSpace:  splitUrl[len(splitUrl)-2],
		Domain:     splitUrl[len(splitUrl)-3],
	}
}

func (r *RemoteInfo) RepositoryFullName() string {
	return fmt.Sprintf("%s/%s", r.NameSpace, r.Repository)
}

func (r *RemoteInfo) BaseUrl() string {
	return "https://" + r.Domain
}

func (r *RemoteInfo) RepositoryUrl() string {
	return strings.Join([]string{r.BaseUrl(), r.RepositoryFullName()}, "/")
}

func (r *RemoteInfo) IssueUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "issues"}, "/")
}

func (r *RemoteInfo) IssueDetailUrl(issueNo int) string {
	return strings.Join([]string{r.IssueUrl(), fmt.Sprintf("%d", issueNo)}, "/")
}

func (r *RemoteInfo) MergeRequestUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "merge_requests"}, "/")
}

func (r *RemoteInfo) MergeRequestDetailUrl(mergeRequestNo int) string {
	return strings.Join([]string{r.MergeRequestUrl(), fmt.Sprintf("%d", mergeRequestNo)}, "/")
}

func (r *RemoteInfo) ApiUrl() string {
	return strings.Join([]string{r.BaseUrl(), "api", "v4"}, "/")
}

func GitRemotes() ([]RemoteInfo, error) {
	// Get remote repositorys
	remotes := utils.GitOutputs("git", []string{"remote"})
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}
	// Extract domain, namespace, repository name from git remote url
	var remoteInfos []RemoteInfo
	for _, remote := range remotes {
		url := utils.GitOutput("git", []string{"remote", "get-url", remote})
		remoteInfo := NewRemoteInfo(url)
		remoteInfos = append(remoteInfos, *remoteInfo)
	}
	return remoteInfos, nil
}
