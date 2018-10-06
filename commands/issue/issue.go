package issue

import (
	"bytes"
	"fmt"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

const IssueTemplateDir = ".gitlab/issue_templates"

type CreateUpdateIssueOption struct {
	Edit       bool   `short:"e" long:"edit" description:"Edit the issue on editor. Start the editor with the contents in the given title and message options."`
	Title      string `short:"i" long:"title" value-name:"<title>" description:"The title of an issue"`
	Message    string `short:"m" long:"message" value-name:"<message>" description:"The message of an issue"`
	Template   string `short:"p" long:"template" value-name:"<issue template>" description:"The template of an issue"`
	StateEvent string `long:"state-event" description:"Change the status. \"close\", \"reopen\""`
	AssigneeID int    `long:"assignee-id" description:"The ID of assignee."`
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
  lab issue [-e] [-i <title>] [-m <message>] [--assignee-id=<assignee id>]

  # Update issue
  lab issue <Issue IID> [-e] [-i <title>] [-m <message>] [--state-event=<state>] [--assignee-id=<assignee id>]

  # Show issue
  lab issue <Issue IID>`
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

	client, err := c.Provider.GetIssueClient(gitlabRemote)
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
		output, err := createIssue(
			client,
			gitlabRemote.RepositoryFullName(),
			createUpdateOption,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case CreateIssueOnEditor:
		var template string
		templateFilename := issueCommandOption.CreateUpdateOption.Template
		if templateFilename != "" {
			res, err := c.getIssueTemplateContent(templateFilename, gitlabRemote)
			if err != nil {
				c.Ui.Error(err.Error())
				return ExitCodeError
			}
			template = res
		}

		output, err := createIssueOnEditor(
			client,
			gitlabRemote.RepositoryFullName(),
			template,
			issueCommandOption.CreateUpdateOption,
			c.EditFunc,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case ShowIssue:
		res, err := detail(client, gitlabRemote.RepositoryFullName(), iid)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case ListIssue:
		listOption := issueCommandOption.ListOption
		res, err := listOfProject(
			client,
			gitlabRemote.RepositoryFullName(),
			listOption,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case ListIssueAllProject:
		// Do get issue list
		listOption := issueCommandOption.ListOption
		res, err := listAll(client, listOption)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

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
		if hasEditIssueOption(createUpdateOption) {
			return UpdateIssue
		}
		return ShowIssue
	}

	// Case of nothing Issue IID
	if createUpdateOption.Edit {
		return CreateIssueOnEditor
	}
	if hasEditIssueOption(createUpdateOption) {
		return CreateIssue
	}
	if listOption.AllProject {
		return ListIssueAllProject
	}

	return ListIssue
}

func hasEditIssueOption(opt *CreateUpdateIssueOption) bool {
	if opt.Title != "" || opt.Message != "" || opt.StateEvent != "" || opt.AssigneeID != 0 {
		return true
	}
	return false
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

func editIssueMessage(title, description string) string {
	message := `%s

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

func (c *IssueCommand) getIssueTemplateContent(templateFilename string, gitlabRemote *git.RemoteInfo) (string, error) {
	issueTemplateClient, err := c.Provider.GetRepositoryClient(gitlabRemote)
	if err != nil {
		return "", err
	}

	filename := IssueTemplateDir + "/" + templateFilename
	res, err := issueTemplateClient.GetFile(
		gitlabRemote.RepositoryFullName(),
		filename,
		makeShowIssueTemplateOption(),
	)
	if err != nil {
		return "", err
	}

	return res, nil
}

func makeShowIssueTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
	}
	return opt
}
