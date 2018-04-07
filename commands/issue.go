package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type CreateUpdateIssueOption struct {
	Edit       bool   `short:"e" long:"edit" description:"Edit the issue on editor. Start the editor with the contents in the given title and message options."`
	Title      string `short:"i" long:"title" value-name:"<title>" description:"The title of an issue"`
	Message    string `short:"m" long:"message" value-name:"<message>" description:"The message of an issue"`
	StateEvent string `long:"state-event" description:"Change the status. \"opened\", \"closed\""`
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
	AssignedMe bool   `short:"a" long:"assigned-me" description:"Shorthand of the scope option for \"--scope=assigned-by-me\"."`
	AllProject bool   `short:"A" long:"all-project" description:"Print the issue of all projects"`
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

type IssueOperation int

const (
	CreateIssue IssueOperation = iota
	CreateIssueOnEditor
	UpdateIssue
	UpdateIssueOnEditor
	ShowIssue
	ListIssue
	ListIssueAllProject
)

func newListIssueOption() *ListIssueOption {
	return &ListIssueOption{}
}

type IssueCommnadOption struct {
	CreateUpdateOption *CreateUpdateIssueOption `group:"Create, Update Options"`
	ListOption         *ListIssueOption         `group:"List Options"`
}

func newIssueOptionParser(opt *IssueCommnadOption) *flags.Parser {
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
  lab issue <Issue IID> [-e] [-i <title>] [-m <message>] 

  # Show issue
  lab issue <Issue IID>
`
	return parser
}

type IssueCommand struct {
	Ui       ui.Ui
	Provider lab.Provider
	EditFunc func(program, file string) error
}

func (c *IssueCommand) Synopsis() string {
	return "Create and Edit, list a issue"
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

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	iid, err := validIssueIID(parseArgs)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Do issue operation
	switch issueOperation(issueCommandOption, parseArgs) {
	case UpdateIssue:
		output, err := updateIssue(
			client,
			gitlabRemote.RepositoryFullName(),
			iid,
			issueCommandOption.CreateUpdateOption,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case UpdateIssueOnEditor:
		output, err := updateIssueOnEditor(
			client,
			gitlabRemote.RepositoryFullName(),
			iid,
			issueCommandOption.CreateUpdateOption,
			c.EditFunc,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case CreateIssue:
		// Do create issue
		createUpdateOption := issueCommandOption.CreateUpdateOption
		issue, err := client.CreateIssue(
			makeCreateIssueOptions(createUpdateOption.Title, createUpdateOption.Message),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print created Issue IID
		c.Ui.Message(fmt.Sprintf("#%d", issue.IID))

	case CreateIssueOnEditor:
		// Starting editor for edit title and description
		createUpdateOption := issueCommandOption.CreateUpdateOption
		template := editIssueMessage(createUpdateOption.Title, createUpdateOption.Message)
		title, message, err := editIssueTitleAndDesc(template, c.EditFunc)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Do create issue
		issue, err := client.CreateIssue(
			makeCreateIssueOptions(title, message),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print created Issue IID
		c.Ui.Message(fmt.Sprintf("#%d", issue.IID))

	case ShowIssue:
		// Do get issue detail
		issue, err := client.GetIssue(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print issue detail
		c.Ui.Message(issueDetailOutput(issue))

	case ListIssue:
		listOption := issueCommandOption.ListOption
		issues, err := client.ProjectIssues(
			makeProjectIssueOption(listOption),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print issue list
		output := projectIssueOutput(issues)
		result := columnize.SimpleFormat(output)
		c.Ui.Message(result)

	case ListIssueAllProject:
		// Do get issue list
		listOption := issueCommandOption.ListOption
		issues, err := client.Issues(makeIssueOption(listOption))
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print issue list
		output := issueOutput(issues)
		result := columnize.SimpleFormat(output)
		c.Ui.Message(result)

	default:
		c.Ui.Error("Invalid issue operation")
		return ExitCodeError
	}

	return ExitCodeOK
}

func issueOperation(opt IssueCommnadOption, args []string) IssueOperation {
	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption

	// Case of getting Issue IID
	if len(args) > 0 {
		if createUpdateOption.Edit {
			return UpdateIssueOnEditor
		}
		if createUpdateOption.Title != "" || createUpdateOption.Message != "" || createUpdateOption.StateEvent != "" {
			return UpdateIssue
		}
		return ShowIssue
	}

	// Case of nothing Issue IID
	if createUpdateOption.Edit {
		return CreateIssueOnEditor
	}
	if createUpdateOption.Title != "" {
		return CreateIssue
	}
	if listOption.AllProject {
		return ListIssueAllProject
	}

	return ListIssue
}

func validIssueIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid Issue IID. IID: %s", args[0])
	}
	return iid, nil
}

func updateIssue(client lab.Client, project string, iid int, opt *CreateUpdateIssueOption) (string, error) {
	// Getting exist issue
	issue, err := client.GetIssue(iid, project)
	if err != nil {
		return "", err
	}

	// Create new title or description
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if opt.Title != "" {
		updatedTitle = opt.Title
	}
	if opt.Message != "" {
		updatedMessage = opt.Message
	}

	// Do update issue
	updatedIssue, err := client.UpdateIssue(
		makeUpdateIssueOption(opt, updatedTitle, updatedMessage),
		iid,
		project,
	)
	if err != nil {
		return "", err
	}

	// Print update Issue IID
	return fmt.Sprintf("#%d", updatedIssue.IID), nil
}

func updateIssueOnEditor(client lab.Client, project string, iid int, opt *CreateUpdateIssueOption, editFunc func(program, file string) error) (string, error) {
	// Getting exist issue
	issue, err := client.GetIssue(iid, project)
	if err != nil {
		return "", err
	}

	// Create new title or description
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if opt.Title != "" {
		updatedTitle = opt.Title
	}
	if opt.Message != "" {
		updatedMessage = opt.Message
	}

	// Starting editor for edit title and description
	template := editIssueMessage(updatedTitle, updatedMessage)
	title, message, err := editIssueTitleAndDesc(template, editFunc)
	if err != nil {
		return "", err
	}

	// Do update issue
	updatedIssue, err := client.UpdateIssue(
		makeUpdateIssueOption(opt, title, message),
		iid,
		project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("#%d", updatedIssue.IID), nil
}

func makeProjectIssueOption(issueListOption *ListIssueOption) *gitlab.ListProjectIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		State:       gitlab.String(issueListOption.GetState()),
		Scope:       gitlab.String(issueListOption.GetScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func makeIssueOption(issueListOption *ListIssueOption) *gitlab.ListIssuesOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: issueListOption.Num,
	}
	listIssuesOptions := &gitlab.ListIssuesOptions{
		State:       gitlab.String(issueListOption.GetState()),
		Scope:       gitlab.String(issueListOption.GetScope()),
		OrderBy:     gitlab.String(issueListOption.OrderBy),
		Sort:        gitlab.String(issueListOption.Sort),
		ListOptions: *listOption,
	}
	return listIssuesOptions
}

func makeCreateIssueOptions(title, description string) *gitlab.CreateIssueOptions {
	opt := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(description),
	}
	return opt
}

func makeUpdateIssueOption(opt *CreateUpdateIssueOption, title, description string) *gitlab.UpdateIssueOptions {
	updateIssueOption := &gitlab.UpdateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(description),
		StateEvent:  gitlab.String(opt.StateEvent),
	}
	return updateIssueOption
}

func issueOutput(issues []*gitlab.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			lab.ParceRepositoryFullName(issue.WebURL),
			fmt.Sprintf("#%d", issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}

func issueDetailOutput(issue *gitlab.Issue) string {
	base := `#%d
Title: %s
Assignee: %s
Author: %s
State: %s
CreatedAt: %s
UpdatedAt: %s

%s`
	detial := fmt.Sprintf(
		base,
		issue.IID,
		issue.Title,
		issue.Assignee.Name,
		issue.Author.Name,
		issue.State,
		issue.CreatedAt.String(),
		issue.UpdatedAt.String(),
		issue.Description,
	)
	return detial
}

func projectIssueOutput(issues []*gitlab.Issue) []string {
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

func editIssueMessage(title, description string) string {
	message := `<!-- Write a message for this issue. The first block of text is the title -->
%s

<!-- the rest is the description.  -->
%s
`
	message = fmt.Sprintf(message, title, description)
	return message
}

func editIssueTitleAndDesc(template string, editFunc func(program, file string) error) (string, string, error) {
	editor, err := git.NewEditor("ISSUE", "issue", template, editFunc)
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

	return title, description, nil
}
