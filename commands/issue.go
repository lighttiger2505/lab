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
	Edit    bool   `short:"e" long:"update" description:"Edit the issue on editor. Start the editor with the contents in the given title and message options."`
	Title   string `short:"i" long:"title" value-name:"<title>" description:"The title of an issue"`
	Message string `short:"m" long:"message" value-name:"<message>" description:"The message of an issue"`
}

func newAddIssueOption() *CreateUpdateIssueOption {
	return &CreateUpdateIssueOption{}
}

type ListIssueOption struct {
	Num        int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of issue to output."`
	State      string `long:"state" value-name:"<state>" default:"all" default-mask:"all" description:"Print only issue of the state just those that are \"opened\", \"closed\" or \"all\""`
	Scope      string `long:"scope" value-name:"<scope>" default:"all" default-mask:"all" description:"Print only given scope. \"created-by-me\", \"assigned-to-me\" or \"all\"."`
	OrderBy    string `long:"orderby" value-name:"<orderby>" default:"updated_at" default-mask:"updated_at" description:"Print issue ordered by \"created_at\" or \"updated_at\" fields."`
	Sort       string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print issue ordered in \"asc\" or \"desc\" order."`
	Opened     bool   `short:"o" long:"opened" description:"Shorthand of the state option for \"--state=opened\"."`
	Closed     bool   `short:"c" long:"closed" description:"Shorthand of the state option for \"--state=closed\"."`
	CreatedMe  bool   `short:"r" long:"created-me" description:"Shorthand of the scope option for \"--scope=created-by-me\"."`
	AssignedMe bool   `long:"s" long:"assigned-me" description:"Shorthand of the scope option for \"--scope=assigned-by-me\"."`
	AllProject bool   `long:"A" long:"all-project" description:"Output the issue of all projects registed in highest priority GitLab server.\n Choise of the GitLab server is determined by \"domains\" setting in \".labconfig.yml\""`
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
	parser.Usage = `issue - Create and Edit, list a issue

Synopsis:
    # List issue
    lab issue [-n <num>] [--state=<state> | -o | -c] [--scope=<scope> | -r | -s]
              [--orderby=<orderby>] [--sort=<sort>] [-A]

    # Create issue
    lab issue [-e] [-i <title>] [-m <message>]

    # Update issue
    lab issue <issue iid> [-e] [-i <title>] [-m <message>] 

    # Show issue
    lab issue <Issue iid>
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
