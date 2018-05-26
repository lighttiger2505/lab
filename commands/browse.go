package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/git"
	gitpath "github.com/lighttiger2505/lab/git/path"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

type BrowseType int

const (
	Issue BrowseType = iota
	MergeRequest
	PipeLine
)

var browseTypePrefix = map[string]BrowseType{
	"#": Issue,
	"i": Issue,
	"I": Issue,
	"!": MergeRequest,
	"m": MergeRequest,
	"M": MergeRequest,
	"p": PipeLine,
	"P": PipeLine,
}

type BrowseCommandOption struct {
	BrowseOption *BrowseOption `group:"Global Options"`
}

func newBrowseOptionParser(opt *BrowseCommandOption) *flags.Parser {
	opt.BrowseOption = newBrowseOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "browse [options]"
	return parser
}

type BrowseCommand struct {
	Ui        ui.Ui
	Provider  lab.Provider
	GitClient git.Client
	Cmd       cmd.Cmd
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse repository page"
}

func (c *BrowseCommand) Help() string {
	buf := &bytes.Buffer{}
	var browseCommnadOption BrowseCommandOption
	browseOptionParser := newBrowseOptionParser(&browseCommnadOption)
	browseOptionParser.WriteHelp(buf)
	return buf.String()
}

func (c *BrowseCommand) Run(args []string) int {
	var browseCommnadOption BrowseCommandOption
	browseOptionParser := newBrowseOptionParser(&browseCommnadOption)
	// Parse option
	parseArgs, err := browseOptionParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Validate option
	browseOption := browseCommnadOption.BrowseOption
	if err := browseOption.IsValid(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	if browseOption.Project != "" {
		namespace, project := browseOption.NameSpaceAndProject()
		gitlabRemote.NameSpace = namespace
		gitlabRemote.Repository = project
	}

	branch, err := c.GitClient.CurrentRemoteBranch(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var url = ""

	if browseOption.Path != "" {
		gitAbsPath, err := gitpath.Abs(browseOption.Path)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		if browseOption.Line != "" {
			url = gitlabRemote.BranchFileWithLine(branch, gitAbsPath, browseOption.Line)
		} else {
			url = gitlabRemote.BranchPath(branch, gitAbsPath)
		}
	} else if browseOption.CurrentPath {
		gitAbsPath, err := gitpath.Current()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		url = gitlabRemote.BranchPath(branch, gitAbsPath)
	} else {
		url, err = getUrlByRemote(gitlabRemote, parseArgs, branch)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
	}

	if err := c.doBrowse(url); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	return ExitCodeOK
}

func getUrlByRemote(gitlabRemote *git.RemoteInfo, args []string, branch string) (string, error) {
	if len(args) > 0 {
		// Gitlab resource page
		browseType, number, err := splitPrefixAndNumber(args[0])
		if err != nil {
			return "", err
		}
		return makeGitlabResourceUrl(gitlabRemote, browseType, number), nil
	} else {
		if branch == "master" {
			// Repository top page
			return gitlabRemote.RepositoryUrl(), nil
		} else {
			// Current branch top page
			return gitlabRemote.BranchUrl(branch), nil
		}
	}
}

func makeGitlabResourceUrl(gitlabRemote *git.RemoteInfo, browseType BrowseType, number int) string {
	var url string
	if number > 0 {
		switch browseType {
		case Issue:
			url = gitlabRemote.IssueDetailUrl(number)
		case MergeRequest:
			url = gitlabRemote.MergeRequestDetailUrl(number)
		case PipeLine:
			url = gitlabRemote.PipeLineDetailUrl(number)
		default:
			url = ""
		}
	} else {
		switch browseType {
		case Issue:
			url = gitlabRemote.IssueUrl()
		case MergeRequest:
			url = gitlabRemote.MergeRequestUrl()
		case PipeLine:
			url = gitlabRemote.PipeLineUrl()
		default:
			url = ""
		}
	}
	return url
}

func splitPrefixAndNumber(arg string) (BrowseType, int, error) {
	for k, v := range browseTypePrefix {
		if strings.HasPrefix(arg, k) {
			numberStr := strings.TrimPrefix(arg, k)
			if numberStr == "" {
				return v, 0, nil
			}

			number, err := strconv.Atoi(numberStr)
			if err != nil {
				return 0, 0, fmt.Errorf("Invalid browse number. \"%s\"", numberStr)
			}
			return v, number, nil
		}
	}
	return 0, 0, fmt.Errorf("Invalid arg. %s", arg)
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}

func (c *BrowseCommand) doBrowse(url string) error {
	browser := searchBrowserLauncher(runtime.GOOS)
	c.Cmd.SetCmd(browser)
	c.Cmd.WithArg(url)
	if err := c.Cmd.Spawn(); err != nil {
		return err
	}
	return nil
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
