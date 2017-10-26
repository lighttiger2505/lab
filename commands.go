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
	prefixArgs := flags.Args()
	if len(prefixArgs) > 0 {
		browseArg, err := NewBrowseArg(prefixArgs[0])
		if err != nil {
			fmt.Println(err.Error())
			return ExitCodeError
		}
		if browseArg.Type == "Issue" {
			cmdOutput(browser, []string{gitlabRemote.IssueDetailUrl(browseArg.No)})
		} else if browseArg.Type == "MergeRequest" {
			cmdOutput(browser, []string{gitlabRemote.MergeRequestDetailUrl(browseArg.No)})
		}
	} else {
		cmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})
	}
	return ExitCodeOK
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
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("project", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
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

	projectId, err := ProjectId(client, gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: 20,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		Scope:       gitlab.String("all"),
		OrderBy:     gitlab.String("created_at"),
		Sort:        gitlab.String("desc"),
		ListOptions: *listOption,
	}
	issues, _, err := client.Issues.ListProjectIssues(projectId, listProjectIssuesOptions)

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
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("browse", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
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

	projectId, err := ProjectId(client, gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: 20,
	}
	listMergeRequestsOptions := &gitlab.ListProjectMergeRequestsOptions{
		Scope:       gitlab.String("all"),
		OrderBy:     gitlab.String("created_at"),
		Sort:        gitlab.String("desc"),
		ListOptions: *listOption,
	}
	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(projectId, listMergeRequestsOptions)
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
