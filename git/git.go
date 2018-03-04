package git

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lighttiger2505/lab/cmd"
)

type Client interface {
	RemoteInfos() ([]*RemoteInfo, error)
	CurrentBranch() (string, error)
}

type GitClient struct {
	Client
}

func NewGitClient() Client {
	return &GitClient{}
}

func (g *GitClient) RemoteInfos() ([]*RemoteInfo, error) {
	return GitRemotes()
}

func (g *GitClient) CurrentBranch() (string, error) {
	// Get remote repositorys
	branches := cmd.GitOutputs("git", []string{"branch"})

	currentPrefix := "*"
	currentBranch := ""
	for _, branch := range branches {
		if strings.HasPrefix(branch, currentPrefix) {
			trimPrefix := strings.TrimPrefix(branch, currentPrefix)
			currentBranch = strings.Trim(trimPrefix, " ")
		}
	}

	if currentBranch == "" {
		return "", errors.New("Not found current branch")
	}
	return currentBranch, nil
}

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

func GitCurrentBranch() (string, error) {
	// Get remote repositorys
	branches := cmd.GitOutputs("git", []string{"branch"})

	currentPrefix := "*"
	currentBranch := ""
	for _, branch := range branches {
		if strings.HasPrefix(branch, currentPrefix) {
			trimPrefix := strings.TrimPrefix(branch, currentPrefix)
			currentBranch = strings.Trim(trimPrefix, " ")
		}
	}

	if currentBranch == "" {
		return "", errors.New("Not found current branch")
	}
	return currentBranch, nil
}

func GitRemotes() ([]*RemoteInfo, error) {
	// Get remote repositorys
	remotes := cmd.GitOutputs("git", []string{"remote"})
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}
	// Extract domain, namespace, repository name from git remote url
	var remoteInfos []*RemoteInfo
	for _, remote := range remotes {
		url := cmd.GitOutput("git", []string{"remote", "get-url", remote})
		remoteInfo := NewRemoteInfo(url)
		remoteInfos = append(remoteInfos, remoteInfo)
	}
	return remoteInfos, nil
}

func GitEditor() (string, error) {
	outputs, err := gitOutput("var", "GIT_EDITOR")
	if err != nil {
		return "", fmt.Errorf("Can't load git var: GIT_EDITOR")
	}
	return os.ExpandEnv(outputs[0]), nil
}

var GlobalFlags []string
var cachedDir string

func GitDir() (string, error) {
	if cachedDir != "" {
		return cachedDir, nil
	}

	outputs, err := gitOutput("rev-parse", "-q", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("Not a git repository (or any of the parent directories): .git")
	}

	var chdir string
	for i, flag := range GlobalFlags {
		if flag == "-C" {
			dir := GlobalFlags[i+1]
			if filepath.IsAbs(dir) {
				chdir = dir
			} else {
				chdir = filepath.Join(chdir, dir)
			}
		}
	}
	gitDir := outputs[0]

	if !filepath.IsAbs(gitDir) {
		if chdir != "" {
			gitDir = filepath.Join(chdir, gitDir)
		}

		gitDir, err := filepath.Abs(gitDir)
		if err != nil {
			return "", err
		}

		gitDir = filepath.Clean(gitDir)
	}

	cachedDir = gitDir
	return gitDir, nil
}

func gitOutput(input ...string) (outputs []string, err error) {
	cmd := gitCmd(input...)

	out, err := cmd.CombinedOutput()
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) != "" {
			outputs = append(outputs, string(line))
		}
	}

	return outputs, err
}

func gitCmd(args ...string) *cmd.BasicCmd {
	cmd := cmd.NewBasicCmd("git")

	for _, v := range GlobalFlags {
		cmd.WithArg(v)
	}

	for _, a := range args {
		cmd.WithArg(a)
	}

	return cmd
}

func cmdOutput(name string, args []string) string {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}

func CommentChar() string {
	char, err := Config("core.commentchar")
	if err != nil {
		char = "#"
	}

	return char
}

func Config(name string) (string, error) {
	return gitGetConfig(name)
}

func gitGetConfig(args ...string) (string, error) {
	output, err := gitOutput(gitConfigCommand(args)...)
	if err != nil {
		return "", fmt.Errorf("Unknown config %s", args[len(args)-1])
	}

	if len(output) == 0 {
		return "", nil
	}

	return output[0], nil
}

func gitConfigCommand(args []string) []string {
	cmd := []string{"config"}
	return append(cmd, args...)
}

type MockClient struct {
	MockRemoteInfos   func() ([]*RemoteInfo, error)
	MockCurrentBranch func() (string, error)
}

func (m *MockClient) RemoteInfos() ([]*RemoteInfo, error) {
	return m.MockRemoteInfos()
}

func (m *MockClient) CurrentBranch() (string, error) {
	return m.MockCurrentBranch()
}
