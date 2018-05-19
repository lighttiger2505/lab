package git

import (
	"fmt"
	"regexp"
	"strings"
)

type RemoteInfo struct {
	Remote     string
	Domain     string
	NameSpace  string
	Repository string
}

func NewRemoteInfo(remote, url string) *RemoteInfo {
	splitUrl := regexp.MustCompile("/|:|@").Split(url, -1)
	return &RemoteInfo{
		Remote:     remote,
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

func (r *RemoteInfo) BranchUrl(branch string) string {
	return strings.Join([]string{r.BaseUrl(), r.RepositoryFullName(), "tree", branch}, "/")
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

func (r *RemoteInfo) PipeLineUrl() string {
	return strings.Join([]string{r.RepositoryUrl(), "pipelines"}, "/")
}

func (r *RemoteInfo) PipeLineDetailUrl(iid int) string {
	return strings.Join([]string{r.PipeLineUrl(), fmt.Sprintf("%d", iid)}, "/")
}

func (r *RemoteInfo) ApiUrl() string {
	return strings.Join([]string{r.BaseUrl(), "api", "v4"}, "/")
}
