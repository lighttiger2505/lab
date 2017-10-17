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

	browser := SearchBrowserLauncher(runtime.GOOS)
	cmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})

	return ExitCodeOK
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

	browser := SearchBrowserLauncher(runtime.GOOS)

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

func SearchBrowserLauncher(goos string) (browser string) {
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
