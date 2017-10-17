package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type ProjectCommand struct {
}

func (c *ProjectCommand) Synopsis() string {
	return "Browse project"
}

func (c *ProjectCommand) Help() string {
	return "Usage: lab project [option]"
}

func (c *ProjectCommand) Run(args []string) int {
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("project", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	gitRemotes, err := GitRemotes()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := FilterGitlabRemote(gitRemotes)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	browser := searchBrowserLauncher(runtime.GOOS)
	cmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})

	return ExitCodeOK
}

type BrowseArg struct {
	Type string
	No   int
}

func NewBrowseArg(arg string) (*BrowseArg, error) {
	var browseArg BrowseArg
	if strings.HasPrefix(arg, "#") {
		number, err := strconv.Atoi(strings.TrimPrefix(arg, "#"))
		if err != nil {
			return nil, errors.New("Invalid number")
		}
		browseArg = BrowseArg{
			Type: "MergeRequest",
			No:   number,
		}
	} else {
		return nil, errors.New("Invalid args")
	}
	return &browseArg, nil
}

type IssueCommand struct {
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	return "Usage: lab issue [option]"
}

func (c *IssueCommand) Run(args []string) int {
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("browse", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	gitRemotes, err := GitRemotes()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := FilterGitlabRemote(gitRemotes)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	browser := searchBrowserLauncher(runtime.GOOS)

	if len(flags.Args()) > 0 {
		issueNo, err := strconv.Atoi(flags.Args()[0])
		if err != nil {
			fmt.Println(err.Error())
		}
		cmdOutput(browser, []string{gitlabRemote.IssueDetailUrl(issueNo)})
	} else {
		cmdOutput(browser, []string{gitlabRemote.IssueUrl()})
	}

	return ExitCodeOK
}

func GitRemotes() ([]GitRemote, error) {
	// Get remote repositorys
	remotes := gitOutputs("git", []string{"remote"})

	// Remote repository is not registered
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}

	var gitRemotes []GitRemote
	for _, remote := range remotes {
		url := gitOutput("git", []string{"remote", "get-url", remote})

		gitRemote, err := NewRemoteUrl(url)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed serialize remote url. %s", url))
		}

		gitRemotes = append(gitRemotes, *gitRemote)
	}

	return gitRemotes, nil
}

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

func searchBrowserLauncher(goos string) (browser string) {
	switch goos {
	case "darwin":
		browser = "open"
	case "windows":
		browser = "cmd /c start"
	default:
		candidates := []string{
			"xdg-open",
			"cygstart",
			"x-www-browser",
			"firefox",
			"opera",
			"mozilla",
			"netscape",
		}
		for _, b := range candidates {
			path, err := exec.LookPath(b)
			if err == nil {
				browser = path
				break
			}
		}
	}
	return browser
}

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"project": func() (cli.Command, error) {
			return &ProjectCommand{}, nil
		},
		"issue": func() (cli.Command, error) {
			return &IssueCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
