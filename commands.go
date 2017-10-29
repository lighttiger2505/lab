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
		fmt.Println(err.Error())
		return ExitCodeError
	}

	gitlabRemoteInfo, err := FilterGitlabRemote(remoteInfos)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	browser := SearchBrowserLauncher(runtime.GOOS)
	prefixArgs := flags.Args()
	if len(prefixArgs) > 0 {
		browseType, number, err := splitPrefixAndNumber(prefixArgs[0])
		if err != nil {
			fmt.Println(err.Error())
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

func (c *IssueCommand) Run(args []string) int {
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
	flags.IntVar(&line, "n", lineDefault, lineHelp)
	flags.IntVar(&line, "line", lineDefault, lineHelp)
	flags.StringVar(&state, "state", "all", "just those that are opened or closed")
	flags.StringVar(&scope, "scope", "all", "given scope: created-by-me, assigned-to-me or all. Defaults to all")
	flags.StringVar(&orderBy, "orderby", "created_at", "ordered by created_at or updated_at fields. Default is created_at")
	flags.StringVar(&sort, "sort", "desc", "sorted in asc or desc order. Default is desc")
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
		PerPage: line,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(state),
		Scope:       gitlab.String(scope),
		OrderBy:     gitlab.String(orderBy),
		Sort:        gitlab.String(sort),
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
	Ui cli.Ui
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	return "Usage: lab merge-request [option]"
}

func (c *MergeRequestCommand) Run(args []string) int {
	var (
		line    int
		state   string
		scope   string
		orderBy string
		sort    string
	)

	// Set subcommand flags
	flags := flag.NewFlagSet("merge-request", flag.ContinueOnError)

	lineHelp := "output the NUM lines"
	lineDefault := 20
	flags.IntVar(&line, "n", lineDefault, lineHelp)
	flags.IntVar(&line, "line", lineDefault, lineHelp)
	flags.StringVar(&state, "state", "all", "just those that are opened or closed")
	flags.StringVar(&scope, "scope", "all", "given scope: created-by-me, assigned-to-me or all. Defaults to all")
	flags.StringVar(&orderBy, "orderby", "created_at", "ordered by created_at or updated_at fields. Default is created_at")
	flags.StringVar(&sort, "sort", "desc", "sorted in asc or desc order. Default is desc")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
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
		PerPage: line,
	}
	listMergeRequestsOptions := &gitlab.ListProjectMergeRequestsOptions{
		State:       gitlab.String(state),
		Scope:       gitlab.String(scope),
		OrderBy:     gitlab.String(orderBy),
		Sort:        gitlab.String(sort),
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
