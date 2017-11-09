package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type BrowseCommand struct {
	Ui cli.Ui
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

	remoteInfos, err := GitRemotes()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemoteInfo, err := FilterGitlabRemote(remoteInfos)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	browser := SearchBrowserLauncher(runtime.GOOS)
	prefixArgs := flags.Args()
	if len(prefixArgs) > 0 {
		browseType, number, err := splitPrefixAndNumber(prefixArgs[0])
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		cmdOutput(browser, []string{browseUrl(gitlabRemoteInfo, browseType, number)})
	} else {
		cmdOutput(browser, []string{gitlabRemoteInfo.RepositoryUrl()})
	}
	return ExitCodeOK
}

func browseUrl(gitlabRemote *RemoteInfo, browseType BrowseType, number int) string {
	var url string
	switch browseType {
	case Issue:
		url = gitlabRemote.IssueDetailUrl(number)
	case MergeRequest:
		url = gitlabRemote.MergeRequestDetailUrl(number)
	default:
		url = ""
	}
	return url
}

type IssueCommand struct {
	Ui cli.Ui
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	return "Usage: lab issue [option]"
}

func overrideArgs(args []string, orArgs []string) []string {
	for _, orArg := range orArgs {
		orArgKey := strings.Split(orArg, "=")[0]
		exist := false
		for _, arg := range args {
			argKey := strings.Split(arg, "=")[0]
			if orArgKey == argKey {
				exist = true
			}
		}
		if exist {
			args = append(args, orArg)
		}
	}
	return args
}

type SearchOpts struct {
	Line    int    `short:"n" long:"line"  discription:"output the NUM lines (default:20)"`
	State   string `short:"t" long:"state" discription:"just those that are opened or closed (default:all)"`
	Scope   string `short:"c" long:"scope" discription:"given scope: created-by-me, assigned-to-me or all. Defaults to all (default:all)"`
	OrderBy string `short:"o" long:"orderby" discription:"ordered by created_at or updated_at fields. Default is created_at (default:created_at)"`
	Sort    string `short:"s" long:"sort" discription:"sorted in asc or desc order. Default is desc (default:desc)"`
}

func NewSearchOpts(args []string, config *Config) (*SearchOpts, error) {

	var searchOpts *SearchOpts

	defaultArgs := []string{
		fmt.Sprintf("--line=%d", 20),
		fmt.Sprintf("--state=%s", "all"),
		fmt.Sprintf("--scope=%s", "all"),
		fmt.Sprintf("--orderby=%s", "created_at"),
		fmt.Sprintf("--sort=%s", "desc"),
	}

	configArgs := []string{
		fmt.Sprintf("--line=%d", config.Line),
		fmt.Sprintf("--state=%s", config.State),
		fmt.Sprintf("--scope=%s", config.Scope),
		fmt.Sprintf("--orderby=%s", config.Orderby),
		fmt.Sprintf("--sort=%s", config.Sort),
	}

	overrideArgs := overrideArgs(defaultArgs, configArgs)
	overrideArgs, err := flags.ParseArgs(&searchOpts, overrideArgs)
	if err != nil {
		return nil, fmt.Errorf("Failed parse default args. %v", overrideArgs)
	}

	// Parse command line options
	args, err = flags.ParseArgs(&searchOpts, args)
	if err != nil {
		return nil, fmt.Errorf("Failed parse args. %v", args)
	}

	return searchOpts, nil
}

func (c *IssueCommand) Run(args []string) int {
	config, err := NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	searchOpts, err := NewSearchOpts(args, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := GitlabRemote()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := GitlabClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: searchOpts.Line,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(searchOpts.State),
		Scope:       gitlab.String(searchOpts.Scope),
		OrderBy:     gitlab.String(searchOpts.OrderBy),
		Sort:        gitlab.String(searchOpts.Sort),
		ListOptions: *listOption,
	}
	// issues, _, err := client.Issues.ListProjectIssues(projectId, listProjectIssuesOptions)
	issues, _, err := client.Issues.ListProjectIssues(
		gitlabRemote.RepositoryFullName(),
		listProjectIssuesOptions,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	for _, issue := range issues {
		data := fmt.Sprintf("#%d", issue.IID) + "|" + issue.Title
		datas = append(datas, data)
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Info(result)
	return ExitCodeOK
}

type MergeRequestCommand struct {
	Ui     cli.Ui
	Config *Config
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	return "Usage: lab merge-request [option]"
}

func (c *MergeRequestCommand) Run(args []string) int {
	config, err := NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	searchOpts, err := NewSearchOpts(args, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := GitlabRemote()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := GitlabClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: searchOpts.Line,
	}
	listMergeRequestsOptions := &gitlab.ListProjectMergeRequestsOptions{
		State:       gitlab.String(searchOpts.State),
		Scope:       gitlab.String(searchOpts.Scope),
		OrderBy:     gitlab.String(searchOpts.OrderBy),
		Sort:        gitlab.String(searchOpts.Sort),
		ListOptions: *listOption,
	}
	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(
		gitlabRemote.RepositoryFullName(),
		listMergeRequestsOptions,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	for _, mergeRequest := range mergeRequests {
		data := fmt.Sprintf("!%d", mergeRequest.IID) + "|" + mergeRequest.Title
		datas = append(datas, data)
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Info(result)

	return ExitCodeOK
}
