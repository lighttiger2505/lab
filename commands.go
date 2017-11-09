package main

import (
	"flag"
	"fmt"
	"runtime"

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

type SearchFlags struct {
	Line    int
	State   string
	Scope   string
	OrderBy string
	Sort    string
}

func NewSearchFlags(config *Config) (*flag.FlagSet, *SearchFlags) {
	var (
		line    int
		state   string
		scope   string
		orderBy string
		sort    string
	)

	// Set subcommand flags
	flags := flag.NewFlagSet("issue", flag.ContinueOnError)

	lineHelp := "output the NUM lines"
	lineDefault := 20
	if config.Line > 0 {
		lineDefault = config.Line
	}
	flags.IntVar(&line, "n", lineDefault, lineHelp)
	flags.IntVar(&line, "line", lineDefault, lineHelp)

	stateDefalut := "all"
	if config.State != "" {
		stateDefalut = config.State
	}
	flags.StringVar(&state, "state", stateDefalut, "just those that are opened or closed")

	scopeDefalut := "all"
	if config.Scope != "" {
		scopeDefalut = config.Scope
	}
	flags.StringVar(&scope, "scope", scopeDefalut, "given scope: created-by-me, assigned-to-me or all. Defaults to all")

	orderbyDefalut := "created_at"
	if config.Orderby != "" {
		orderbyDefalut = config.Orderby
	}
	flags.StringVar(&orderBy, "orderby", orderbyDefalut, "ordered by created_at or updated_at fields. Default is created_at")

	sortDefalut := "desc"
	if config.Sort != "" {
		sortDefalut = config.Sort
	}
	flags.StringVar(&sort, "sort", sortDefalut, "sorted in asc or desc order. Default is desc")

	searchFlags := SearchFlags{
		Line:    line,
		State:   state,
		Scope:   scope,
		OrderBy: orderBy,
		Sort:    sort,
	}

	return flags, &searchFlags
}

func (c *IssueCommand) Run(args []string) int {
	config, err := NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	flags, option := NewSearchFlags(config)

	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
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
		PerPage: option.Line,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(option.State),
		Scope:       gitlab.String(option.Scope),
		OrderBy:     gitlab.String(option.OrderBy),
		Sort:        gitlab.String(option.Sort),
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

	flags, option := NewSearchFlags(config)

	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
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
		PerPage: option.Line,
	}
	listMergeRequestsOptions := &gitlab.ListProjectMergeRequestsOptions{
		State:       gitlab.String(option.State),
		Scope:       gitlab.String(option.Scope),
		OrderBy:     gitlab.String(option.OrderBy),
		Sort:        gitlab.String(option.Sort),
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
