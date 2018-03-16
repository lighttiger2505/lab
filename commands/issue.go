package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type CreateUpdateIssueOption struct {
	Edit    bool   `short:"e" long:"update" description:"edit issue"`
	Title   string `short:"i" long:"title" description:"issue issue"`
	Message string `short:"m" long:"message" description:"issue message"`
}

func newAddIssueOption() *CreateUpdateIssueOption {
	return &CreateUpdateIssueOption{}
}

type ListIssueOption struct {
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

func (l *ListIssueOption) GetState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	return l.State
}

func (l *ListIssueOption) GetScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

func newListIssueOption() *ListIssueOption {
	return &ListIssueOption{}
}

type IssueCommnadOption struct {
	GlobalOption       *GlobalOption            `group:"Global Options"`
	CreateUpdateOption *CreateUpdateIssueOption `group:"Create, Update Options"`
	ListOption         *ListIssueOption         `group:"List Options"`
}

func newIssueOptionParser(opt *IssueCommnadOption) *flags.Parser {
	opt.GlobalOption = newGlobalOption()
	opt.CreateUpdateOption = newAddIssueOption()
	opt.ListOption = newListIssueOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = `add-issue [options]

Synopsis:
    lab issue -a <title> [-d <message>]
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
	parseArgs, err := issueCommnadParser.ParseArgs(args)
	if err != nil {
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

	createUpdateOption := issueCommandOption.CreateUpdateOption
	if len(parseArgs) > 0 {
		iid, err := strconv.Atoi(parseArgs[0])
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Invalid issue iid. \"%s\"", parseArgs[0]))
			return ExitCodeError
		}

		issue, err := client.GetIssue(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		if createUpdateOption.Edit {
			title, message, err := editIssueTitleAndDesc(issue.Title, issue.Description)
			if err != nil {
				c.Ui.Error(err.Error())
				return ExitCodeError
			}

			issue, err := client.UpdateIssue(
				makeUpdateIssueOption(title, message),
				iid,
				gitlabRemote.RepositoryFullName(),
			)
			if err != nil {
				c.Ui.Error(err.Error())
				return ExitCodeError
			}
			c.Ui.Message(fmt.Sprintf("#%d", issue.IID))
			return ExitCodeOK
		} else {
			c.Ui.Message(issueDetailOutput(issue))
			return ExitCodeOK
		}
	}

	if createUpdateOption.Edit {
		argsTitle := ""
		if len(parseArgs) > 0 {
			argsTitle = parseArgs[0]
		}

		title, message, err := getIssueTitleAndDesc(argsTitle, createUpdateOption.Message)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		issue, err := client.CreateIssue(
			makeCreateIssueOptions(title, message),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(fmt.Sprintf("#%d", issue.IID))
		return ExitCodeOK
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

func makeProjectIssueOption(issueListOption *ListIssueOption) *gitlabc.ListProjectIssuesOptions {
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

func makeIssueOption(issueListOption *ListIssueOption) *gitlabc.ListIssuesOptions {
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

func makeCreateIssueOptions(title, description string) *gitlabc.CreateIssueOptions {
	opt := &gitlabc.CreateIssueOptions{
		Title:       gitlabc.String(title),
		Description: gitlabc.String(description),
	}
	return opt
}

func makeUpdateIssueOption(title, description string) *gitlabc.UpdateIssueOptions {
	opt := &gitlabc.UpdateIssueOptions{
		Title:       gitlabc.String(title),
		Description: gitlabc.String(description),
	}
	return opt
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

func issueDetailOutput(issue *gitlabc.Issue) string {
	base := `#%d
Title: %s
Assignee: %s
Author: %s
CreatedAt: %s
UpdatedAt: %s

%s`
	detial := fmt.Sprintf(
		base,
		issue.IID,
		issue.Title,
		issue.Assignee.Name,
		issue.Author.Name,
		issue.CreatedAt.String(),
		issue.UpdatedAt.String(),
		issue.Description,
	)
	return detial
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

func createIssueMessage(title, description string) string {
	message := `<!-- Write a message for this issue. The first block of text is the title -->
%s

<!-- the rest is the description.  -->
%s
`
	message = fmt.Sprintf(message, title, description)
	return message
}

func getIssueTitleAndDesc(titleIn, descIn string) (string, string, error) {
	var title, description string
	if titleIn == "" || descIn == "" {
		message := createIssueMessage(titleIn, descIn)

		editor, err := git.NewEditor("ISSUE", "issue", message)
		if err != nil {
			return "", "", err
		}

		title, description, err = editor.EditTitleAndDescription()
		if err != nil {
			return "", "", err
		}

		if editor != nil {
			defer editor.DeleteFile()
		}
	} else {
		title = titleIn
		description = descIn
	}

	if title == "" {
		return "", "", fmt.Errorf("Title is requeired")
	}

	return title, description, nil
}

func editIssueTitleAndDesc(titleIn, descIn string) (string, string, error) {
	message := createIssueMessage(titleIn, descIn)
	editor, err := git.NewEditor("ISSUE", "issue", message)
	if err != nil {
		return "", "", err
	}

	title, description, err := editor.EditTitleAndDescription()
	if err != nil {
		return "", "", err
	}

	if editor != nil {
		defer editor.DeleteFile()
	}

	if title == "" {
		return "", "", fmt.Errorf("Title is requeired")
	}

	return title, description, nil
}
