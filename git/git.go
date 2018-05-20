package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Client interface {
	RemoteInfos() ([]*RemoteInfo, error)
	CurrentBranch(remote *RemoteInfo) (string, error)
}

type GitClient struct {
	Client
}

func NewGitClient() Client {
	return &GitClient{}
}

func (g *GitClient) RemoteInfos() ([]*RemoteInfo, error) {
	// Get remote repositorys
	remotes, err := gitOutput("remote")
	if err != nil {
		return nil, fmt.Errorf("Failed collect git remove infos. %s", err)
	}
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}

	// Extract domain, namespace, repository name from git remote url
	var remoteInfos []*RemoteInfo
	for _, remote := range remotes {
		url, err := gitOutput("remote", "get-url", remote)
		fmt.Println(url)
		if err != nil {
			return nil, fmt.Errorf("Failed get git remote url. %s", err)
		}
		remoteInfo := NewRemoteInfo(remote, url[0])
		remoteInfos = append(remoteInfos, remoteInfo)
	}
	return remoteInfos, nil
}

func (g *GitClient) CurrentBranch(remote *RemoteInfo) (string, error) {
	// Get remote repositorys
	branches, err := gitOutput("branch", "-a")
	if err != nil {
		return "", fmt.Errorf("Failed get git branches. %s", err)
	}

	currentPrefix := "*"
	currentBranch := ""
	for _, branch := range branches {
		if strings.HasPrefix(branch, currentPrefix) {
			trimPrefix := strings.TrimPrefix(branch, currentPrefix)
			currentBranch = strings.Trim(trimPrefix, " ")
			break
		}
	}

	if currentBranch == "" {
		return "", errors.New("Not found current branch")
	}

	remoteBranch := fmt.Sprintf("%s/%s", remote.Remote, currentBranch)
	for _, branch := range branches {
		trimBranch := strings.TrimSpace(branch)
		if strings.HasSuffix(trimBranch, remoteBranch) {
			return currentBranch, nil
		}
	}
	return "master", nil

}

func GitCurrentBranch() (string, error) {
	// Get remote repositorys
	branches, err := gitOutput("branch")
	if err != nil {
		return "", fmt.Errorf("Failed get git branches. %s", err)
	}

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

func GitEditor() (string, error) {
	outputs, err := gitOutput("var", "GIT_EDITOR")
	if err != nil {
		return "", fmt.Errorf("Can't load git var: GIT_EDITOR")
	}
	return os.ExpandEnv(outputs[0]), nil
}

var cachedDir string

func GitDir() (string, error) {
	if cachedDir != "" {
		return cachedDir, nil
	}

	outputs, err := gitOutput("rev-parse", "-q", "--git-dir")
	if err != nil {
		return "", fmt.Errorf("Not a git repository (or any of the parent directories): .git")
	}

	gitDir := outputs[0]
	if !filepath.IsAbs(gitDir) {
		gitDir, err := filepath.Abs(gitDir)
		if err != nil {
			return "", err
		}
		gitDir = filepath.Clean(gitDir)
	}

	cachedDir = gitDir
	return gitDir, nil
}

// For os/exec test
var execCommand = exec.Command

func gitOutput(input ...string) (outputs []string, err error) {
	output, err := execCommand("git", input...).CombinedOutput()
	for _, line := range strings.Split(string(output), "\n") {
		if strings.TrimSpace(line) != "" {
			outputs = append(outputs, string(line))
		}
	}
	return outputs, err
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

func (m *MockClient) CurrentBranch(remote *RemoteInfo) (string, error) {
	return m.MockCurrentBranch()
}
