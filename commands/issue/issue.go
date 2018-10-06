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

const TemplateDir = ".gitlab/issue_templates"

type CreateUpdateOption struct {
	Edit       bool   `short:"e" long:"edit" description:"Edit the issue on editor. Start the editor with the contents in the given title and message options."`
	Title      string `short:"i" long:"title" value-name:"<title>" description:"The title of an issue"`
	Message    string `short:"m" long:"message" value-name:"<message>" description:"The message of an issue"`
	Template   string `short:"p" long:"template" value-name:"<issue template>" description:"The template of an issue"`
	StateEvent string `long:"state-event" description:"Change the status. \"close\", \"reopen\""`
	AssigneeID int    `long:"assignee-id" description:"The ID of assignee."`
}

func newCreateUpdateOption() *CreateUpdateOption {
	return &CreateUpdateOption{}
}

type ListOption struct {
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

func (l *ListOption) GetState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	return l.State
}

func (l *ListOption) GetScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

type Operation int

const (
	Create Operation = iota
	CreateOnEditor
	Update
	UpdateOnEditor
	Detail
	ListIssue
	List
)

func newOption() *ListOption {
	return &ListOption{}
}

type Option struct {
	CreateUpdateOption *CreateUpdateOption `group:"Create, Update Options"`
	ListOption         *ListOption         `group:"List Options"`
}

func newOptionParser(opt *Option) *flags.Parser {
	opt.CreateUpdateOption = newCreateUpdateOption()
	opt.ListOption = newOption()
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
	var opt Option
	parser := newOptionParser(&opt)
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	var opt Option
	parser := newOptionParser(&opt)
	parseArgs, err := parser.ParseArgs(args)
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
	switch getOperation(opt, parseArgs) {
	case Update:
		output, err := update(
			client,
			gitlabRemote.RepositoryFullName(),
			iid,
			opt.CreateUpdateOption,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case UpdateOnEditor:
		output, err := updateOnEditor(
			client,
			gitlabRemote.RepositoryFullName(),
			iid,
			opt.CreateUpdateOption,
			c.EditFunc,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case Create:
		// Do create issue
		createUpdateOption := opt.CreateUpdateOption
		output, err := create(
			client,
			gitlabRemote.RepositoryFullName(),
			createUpdateOption,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case CreateOnEditor:
		var template string
		templateFilename := opt.CreateUpdateOption.Template
		if templateFilename != "" {
			res, err := c.getIssueTemplateContent(templateFilename, gitlabRemote)
			if err != nil {
				c.Ui.Error(err.Error())
				return ExitCodeError
			}
			template = res
		}

		output, err := createOnEditor(
			client,
			gitlabRemote.RepositoryFullName(),
			template,
			opt.CreateUpdateOption,
			c.EditFunc,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(output)

	case Detail:
		res, err := detail(client, gitlabRemote.RepositoryFullName(), iid)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case ListIssue:
		listOption := opt.ListOption
		res, err := list(
			client,
			gitlabRemote.RepositoryFullName(),
			listOption,
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case List:
		// Do get issue list
		listOption := opt.ListOption
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

func getOperation(opt Option, args []string) Operation {
	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption

	// Case of getting Issue IID
	if len(args) > 0 {
		if createUpdateOption.Edit {
			return UpdateOnEditor
		}
		if hasEditIssueOption(createUpdateOption) {
			return Update
		}
		return Detail
	}

	// Case of nothing Issue IID
	if createUpdateOption.Edit {
		return CreateOnEditor
	}
	if hasEditIssueOption(createUpdateOption) {
		return Create
	}
	if listOption.AllProject {
		return List
	}

	return ListIssue
}

func hasEditIssueOption(opt *CreateUpdateOption) bool {
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

	filename := TemplateDir + "/" + templateFilename
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
