package git

import (
	"fmt"
	"regexp"
	"strings"
)

type RemoteInfo struct {
	Remote     string
	Domain     string
	Group      string
	SubGroup   string
	Repository string
}

func NewRemoteInfo(remote, url string) *RemoteInfo {
	url = strings.TrimPrefix(url, "ssh://")
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "git@")
	splitUrl := regexp.MustCompile("/|:").Split(url, -1)
	if len(splitUrl) > 3 {
		// Apply subgroup
		return &RemoteInfo{
			Remote:     remote,
			Repository: strings.TrimSuffix(splitUrl[len(splitUrl)-1], ".git"),
			SubGroup:   splitUrl[len(splitUrl)-2],
			Group:      splitUrl[len(splitUrl)-3],
			Domain:     splitUrl[len(splitUrl)-4],
		}
	}
	return &RemoteInfo{
		Remote:     remote,
		Repository: strings.TrimSuffix(splitUrl[len(splitUrl)-1], ".git"),
		Group:      splitUrl[len(splitUrl)-2],
		Domain:     splitUrl[len(splitUrl)-3],
	}
}

func (r *RemoteInfo) RepositoryFullName() string {
	if r.SubGroup != "" {
		return fmt.Sprintf("%s/%s/%s", r.Group, r.SubGroup, r.Repository)
	}
	return fmt.Sprintf("%s/%s", r.Group, r.Repository)
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

func (r *RemoteInfo) BranchPath(branch string, path string) string {
	return strings.Join([]string{r.BranchUrl(branch), path}, "/")
}

func (r *RemoteInfo) BranchFileWithLine(branch string, path string, line string) string {
	return strings.Join([]string{r.BranchPath(branch, path), line}, "/")
}

func (r *RemoteInfo) Subpage(subpage string) string {
	return strings.Join([]string{r.RepositoryUrl(), subpage}, "/")
}

func (r *RemoteInfo) ApiUrl() string {
	return strings.Join([]string{r.BaseUrl(), "api", "v4"}, "/")
}
