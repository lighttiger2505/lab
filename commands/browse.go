package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	url, err := c.getURL(parseArgs, gitlabRemote, branch, browseOption)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	if err := c.doBrowse(url); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	return ExitCodeOK
}

func (c *BrowseCommand) getURL(args []string, remote *git.RemoteInfo, branch string, opt *BrowseOption) (string, error) {
	if len(args) < 1 {
		// TODO You need to ignore the branch when the project is specified as an option
		if branch == "master" {
			return remote.RepositoryUrl(), nil
		}
		return remote.BranchUrl(branch), nil
	}

	arg := args[0]
	// TODO In order to display an appropriate error message, it is necessary to check whether the argument is a file path
	if isFilePath(arg) {
		result, err := isDir(arg)
		if err != nil {
			return "", err
		}
		if result {
			gitAbsPath, err := gitpath.Current()
			if err != nil {
				return "", err
			}
			return remote.BranchPath(branch, gitAbsPath), nil
		}

		gitAbsPath, err := gitpath.Abs(arg)
		if err != nil {
			return "", err
		}
		if opt.Line != "" {
			return remote.BranchFileWithLine(branch, gitAbsPath, opt.Line), nil
		}
		return remote.BranchPath(branch, gitAbsPath), nil
	}

	browseType, number, err := splitPrefixAndNumber(args[0])
	if err != nil {
		return "", err
	}
	return makeGitlabResourceUrl(remote, browseType, number), nil
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

func isFilePath(value string) bool {
	absPath, _ := filepath.Abs(value)
	if isFileExist(absPath) {
		return true
	}
	return false
}

func isDir(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	fi, err := file.Stat()
	switch {
	case err != nil:
		return false, err
	case fi.IsDir():
		return true, nil
	default:
		return false, nil
	}
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
