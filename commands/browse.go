package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
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

var browseOpt BrowseOpt

type BrowseOpt struct {
	GlobalOpt *GlobalOpt `group:"Global Options"`
}

func newBrowseOptionParser(browseOpt *BrowseOpt) *flags.Parser {
	parser := flags.NewParser(issueOpt, flags.Default)
	parser.Usage = "browse [options] [args]"
	return parser
}

type BrowseCommand struct {
	Ui ui.Ui
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse repository page"
}

func (c *BrowseCommand) Help() string {
	buf := &bytes.Buffer{}
	newBrowseOptionParser(&browseOpt).WriteHelp(buf)
	return buf.String()
}

func (c *BrowseCommand) Run(args []string) int {
	parser := newBrowseOptionParser(&browseOpt)
	if _, err := parser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	parseArgs, err := parser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	config, err := config.NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := gitlab.GitlabRemote(c.Ui, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Replace specific repository
	oneDomain := config.PreferredDomains[0]
	if browseOpt.GlobalOpt.Repository != "" {
		namespace, project, err := browseOpt.GlobalOpt.ValidRepository()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		gitlabRemote.Domain = oneDomain
		gitlabRemote.NameSpace = namespace
		gitlabRemote.Repository = project
	}

	// Getting browse command
	browser := searchBrowserLauncher(runtime.GOOS)

	// Browse current repository page
	if gitlabRemote != nil {
		if len(parseArgs) > 0 {
			// Browse github resource page
			browseType, number, err := splitPrefixAndNumber(parseArgs[0])
			if err != nil {
				c.Ui.Error(err.Error())
				return ExitCodeError
			}
			cmd.CmdOutput(browser, []string{browseUrl(gitlabRemote, browseType, number)})
		} else {
			// Browse current branch top page
			currentBranch, err := git.GitCurrentBranch()
			if err != nil {
				c.Ui.Error(err.Error())
				return ExitCodeError
			}
			if currentBranch == "master" || browseOpt.GlobalOpt.Repository != "" {
				cmd.CmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})
			} else {
				cmd.CmdOutput(browser, []string{gitlabRemote.BranchUrl(currentBranch)})
			}
		}
	} else {
		if oneDomain != "" {
			// Browse current domain page
			cmd.CmdOutput(browser, []string{"https://" + oneDomain})
		} else {
			c.Ui.Message("Not found browse url.")
		}
	}

	return ExitCodeOK
}

func browseUrl(gitlabRemote *git.RemoteInfo, browseType BrowseType, number int) string {
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

func splitPrefixAndNumber(arg string) (BrowseType, int, error) {
	for k, v := range browseTypePrefix {
		if strings.HasPrefix(arg, k) {
			numberStr := strings.TrimPrefix(arg, k)
			if numberStr == "" {
				return v, 0, nil
			} else {
				number, err := strconv.Atoi(numberStr)
				if err != nil {
					return 0, 0, errors.New(fmt.Sprintf("Invalid browsing number: %s", arg))
				}
				return v, number, nil
			}
		}
	}
	return 0, 0, errors.New(fmt.Sprintf("Invalid arg: %s", arg))
}
