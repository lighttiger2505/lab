package main

import (
	"bytes"
	"flag"
	"fmt"
	"runtime"

	"github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type SearchOptons struct {
	Line    int    `short:"n" long:"line" default:"20" default-mask:"20" description:"output the NUM lines"`
	State   string `short:"t" long:"state" default:"all" default-mask:"all" description:"just those that are opened, closed or all"`
	Scope   string `short:"c" long:"scope" default:"all" default-mask:"all" description:"given scope: created-by-me, assigned-to-me or all."`
	OrderBy string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by created_at or updated_at fields."`
	Sort    string `short:"s" long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order."`
}

var searchOptions SearchOptons
var searchParser = flags.NewParser(&searchOptions, flags.Default)

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

	config, err := NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := GitlabRemote(c.Ui, config)
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
		cmdOutput(browser, []string{browseUrl(gitlabRemote, browseType, number)})
	} else {
		cmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})
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
	Ui ui.Ui
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	searchParser.Usage = "issue [options]"
	searchParser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	if _, err := searchParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	config, err := NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := GitlabRemote(c.Ui, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := GitlabClient(c.Ui, gitlabRemote, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: searchOptions.Line,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(searchOptions.State),
		Scope:       gitlab.String(searchOptions.Scope),
		OrderBy:     gitlab.String(searchOptions.OrderBy),
		Sort:        gitlab.String(searchOptions.Sort),
		ListOptions: *listOption,
	}

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
	c.Ui.Message(result)
	return ExitCodeOK
}

type MergeRequestCommand struct {
	Ui ui.Ui
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	searchParser.Usage = "merge-request [options]"
	searchParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	if _, err := searchParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	config, err := NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := GitlabRemote(c.Ui, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := GitlabClient(c.Ui, gitlabRemote, config)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: searchOptions.Line,
	}
	listMergeRequestsOptions := &gitlab.ListProjectMergeRequestsOptions{
		State:       gitlab.String(searchOptions.State),
		Scope:       gitlab.String(searchOptions.Scope),
		OrderBy:     gitlab.String(searchOptions.OrderBy),
		Sort:        gitlab.String(searchOptions.Sort),
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
	c.Ui.Message(result)

	return ExitCodeOK
}
