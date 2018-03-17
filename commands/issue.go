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
	AllProject bool   `long:"A" long:"all-project" description:"Print the issue of all projects"`
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
	Create IssueOperation = iota
	CreateOnEditor
	Update
	UpdateOnEditor
	Show
	List
	ListAllProject
)

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

	iid, err := validIID(parseArgs)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Do issue operation
	operation := issueOperation(issueCommandOption, parseArgs)

	switch operation {
	case Update:
		// Getting exist issue
		issue, err := client.GetIssue(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Create new title or description
		createUpdateOption := issueCommandOption.CreateUpdateOption
		updatedTitle := issue.Title
		updatedMessage := issue.Description
		if createUpdateOption.Title != "" {
			updatedTitle = createUpdateOption.Title
		}
		if createUpdateOption.Message != "" {
			updatedMessage = createUpdateOption.Message
		}

		// Do update issue
		updatedIssue, err := client.UpdateIssue(
			makeUpdateIssueOption(updatedTitle, updatedMessage),
			iid,
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print update Issue IID
		c.Ui.Message(fmt.Sprintf("#%d", updatedIssue.IID))

	case UpdateOnEditor:
		// Getting exist issue
		issue, err := client.GetIssue(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Starting editor for edit title and description
		title, message, err := editIssueTitleAndDesc(issue.Title, issue.Description)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Do update issue
		updatedIssue, err := client.UpdateIssue(
			makeUpdateIssueOption(title, message),
			iid,
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(fmt.Sprintf("#%d", updatedIssue.IID))

	case Create:
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

	case CreateOnEditor:
		// Starting editor for edit title and description
		createUpdateOption := issueCommandOption.CreateUpdateOption
		title, message, err := editIssueTitleAndDesc(createUpdateOption.Title, createUpdateOption.Message)
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

	case Show:
		// Do get issue detail
		issue, err := client.GetIssue(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print issue detail
		c.Ui.Message(issueDetailOutput(issue))

	case List:
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

	case ListAllProject:
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
			return UpdateOnEditor
		}
		if createUpdateOption.Title != "" || createUpdateOption.Message != "" {
			return Update
		}
		return Show
	}

	// Case of nothing Issue IID
	if createUpdateOption.Edit {
		return CreateOnEditor
	}
	if createUpdateOption.Title != "" {
		return Create
	}
	if listOption.AllProject {
		return ListAllProject
	}

	return List
}

func validIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid Issue IID. IID: %s", args[0])
	}
	return iid, nil
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

func editIssueTitleAndDesc(title, message string) (string, string, error) {
	template := createIssueMessage(title, message)

	editor, err := git.NewEditor("ISSUE", "issue", template)
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
