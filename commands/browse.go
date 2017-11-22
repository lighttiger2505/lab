package commands

import (
	"flag"
	"runtime"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/lighttiger2505/lab/utils"
)

type BrowseCommand struct {
	Ui ui.Ui
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse project"
}

func (c *BrowseCommand) Help() string {
	return "Usage: lab project [option]"
}

func (c *BrowseCommand) Run(args []string) int {
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("project", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
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

	browser := utils.SearchBrowserLauncher(runtime.GOOS)
	prefixArgs := flags.Args()
	if len(prefixArgs) > 0 {
		browseType, number, err := utils.SplitPrefixAndNumber(prefixArgs[0])
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		utils.CmdOutput(browser, []string{browseUrl(gitlabRemote, browseType, number)})
	} else {
		utils.CmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})
	}
	return ExitCodeOK
}

func browseUrl(gitlabRemote *git.RemoteInfo, browseType utils.BrowseType, number int) string {
	var url string
	switch browseType {
	case utils.Issue:
		url = gitlabRemote.IssueDetailUrl(number)
	case utils.MergeRequest:
		url = gitlabRemote.MergeRequestDetailUrl(number)
	default:
		url = ""
	}
	return url
}
