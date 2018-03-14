package commands

import (
	"bytes"
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type IssueAddOption struct {
	Title       string `short:"t" long:"title" description:"issue title"`
	Description string `short:"d" long:"descript" description:"issue description"`
}

func newIssueAddOption() *IssueAddOption {
	return &IssueAddOption{}
}

type IssueListOption struct {
	List       bool   `short:"l" long:"line" description:"show list"`
	Num        int    `short:"n" long:"num"  default:"20" default-mask:"20" description:"show issue num"`
	State      string `long:"state" default:"all" default-mask:"all" description:"just those that are opened, closed or all"`
	Scope      string `long:"scope" default:"all" default-mask:"all" description:"given scope: created-by-me, assigned-to-me or all."`
	OrderBy    string `long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by created_at or updated_at fields."`
	Sort       string `long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order."`
	Opened     bool   `short:"o" long:"opened" description:"search state opened"`
	Closed     bool   `short:"c" long:"closed" description:"search scope closed"`
	CreatedMe  bool   `short:"r" long:"created-me" description:"search scope created-by-me"`
	AssignedMe bool   `long:"s" long:"assigned-me" description:"search scope assigned-to-me"`
	AllProject bool   `long:"a" long:"all-project" description:"search target all project"`
}

func (l *IssueListOption) GetState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	return l.State
}

func (l *IssueListOption) GetScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

func newIssueListOption() *IssueListOption {
	return &IssueListOption{}
}

type IssueCommnadOption struct {
	GlobalOption *GlobalOption `group:"Global Options"`
	AddOption    *IssueAddOption
	ListOption   *IssueListOption
}

func newIssueOptionParser(opt *IssueCommnadOption) *flags.Parser {
	opt.GlobalOption = newGlobalOption()
	opt.AddOption = newIssueAddOption()
	opt.ListOption = newIssueListOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = `add-issue [options]

Synopsis:
    lab issue [-t <title>] [-d <description>]
    lab issue [-n <num>] -l [--state <state>] [--scope <scope>]
              [--orderby <orderby>] [--sort <sort>] -o -c 
              -cm -am -al
	lab issue [-t <title>] [-d <description>] <id>
`
	return parser
}

type IssueCommand struct {
	Ui       ui.Ui
	Provider gitlab.Provider
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	var issueCommandOption IssueCommnadOption
	issueCommnadParser := newIssueOptionParser(&issueCommandOption)
	issueCommnadParser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	var issueCommandOption IssueCommnadOption
	issueCommnadParser := newIssueOptionParser(&issueCommandOption)
	if _, err := issueCommnadParser.ParseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	globalOption := issueCommandOption.GlobalOption
	if err := globalOption.IsValid(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	var gitlabRemote *git.RemoteInfo
	if globalOption.Project != "" {
		namespace, project := globalOption.NameSpaceAndProject()
		gitlabRemote = c.Provider.GetSpecificRemote(namespace, project)
	} else {
		var err error
		gitlabRemote, err = c.Provider.GetCurrentRemote()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
	}

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	listOption := issueCommandOption.ListOption
	if listOption.AllProject {
		issues, err := client.Issues(makeIssueOption(listOption))
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		datas = issueOutput(issues)

	} else {
		issues, err := client.ProjectIssues(
			makeProjectIssueOption(listOption),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		datas = projectIssueOutput(issues)
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Message(result)

	return ExitCodeOK
}

func makeProjectIssueOption(issueListOption *IssueListOption) *gitlabc.ListProjectIssuesOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listProjectIssuesOptions := &gitlabc.ListProjectIssuesOptions{
		State:       gitlabc.String(issueListOption.GetState()),
		Scope:       gitlabc.String(issueListOption.GetScope()),
		OrderBy:     gitlabc.String(issueListOption.OrderBy),
		Sort:        gitlabc.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func makeIssueOption(issueListOption *IssueListOption) *gitlabc.ListIssuesOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listIssuesOptions := &gitlabc.ListIssuesOptions{
		State:       gitlabc.String(issueListOption.GetState()),
		Scope:       gitlabc.String(issueListOption.GetScope()),
		OrderBy:     gitlabc.String(issueListOption.OrderBy),
		Sort:        gitlabc.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listIssuesOptions
}

func issueOutput(issues []*gitlabc.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("#%d", issue.IID),
			gitlab.ParceRepositoryFullName(issue.WebURL),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}

func projectIssueOutput(issues []*gitlabc.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("#%d", issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}
