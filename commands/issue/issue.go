package issue

import (
	"bytes"
	"fmt"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type CreateUpdateOption struct {
	Edit       bool   `short:"e" long:"edit" description:"Edit the issue on editor. Start the editor with the contents in the given title and message options."`
	Title      string `short:"i" long:"title" value-name:"<title>" description:"The title of an issue"`
	Message    string `short:"m" long:"message" value-name:"<message>" description:"The message of an issue"`
	Template   string `short:"p" long:"template" value-name:"<issue template>" description:"The template of an issue"`
	StateEvent string `long:"state-event" description:"Change the status. \"close\", \"reopen\""`
	AssigneeID int    `long:"assignee-id" description:"The ID of assignee."`
}

func (o *CreateUpdateOption) hasEdit() bool {
	if o.Edit {
		return true
	}
	return false
}

func (o *CreateUpdateOption) hasCreate() bool {
	if o.Title != "" {
		return true
	}
	return false
}

func (o *CreateUpdateOption) hasUpdate() bool {
	if o.Title != "" || o.Message != "" || o.StateEvent != "" || o.AssigneeID != 0 {
		return true
	}
	return false
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

func (l *ListOption) getState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	return l.State
}

func (l *ListOption) getScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

type Option struct {
	CreateUpdateOption *CreateUpdateOption `group:"Create, Update Options"`
	ListOption         *ListOption         `group:"List Options"`
}

func newOptionParser(opt *Option) *flags.Parser {
	opt.CreateUpdateOption = &CreateUpdateOption{}
	opt.ListOption = &ListOption{}
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

	method, err := c.getMethod(opt, parseArgs, gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	res, err := method.Process()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	c.Ui.Message(res)

	return ExitCodeOK
}

func (c *IssueCommand) getMethod(opt Option, args []string, remote *git.RemoteInfo) (internal.Method, error) {
	issueClient, err := c.Provider.GetIssueClient(remote)
	if err != nil {
		return nil, err
	}

	repositoryClient, err := c.Provider.GetRepositoryClient(remote)
	if err != nil {
		return nil, err
	}

	iid, err := validIssueIID(args)
	if err != nil {
		return nil, err
	}

	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption
	project := remote.RepositoryFullName()

	// Case of getting Issue IID
	if len(args) > 0 {
		if createUpdateOption.hasEdit() {
			return &updateOnEditorMethod{
				client:   issueClient,
				opt:      createUpdateOption,
				project:  project,
				id:       iid,
				editFunc: c.EditFunc,
			}, nil
		}
		if createUpdateOption.hasUpdate() {
			return &updateMethod{
				client:  issueClient,
				opt:     createUpdateOption,
				project: project,
				id:      iid,
			}, nil
		}
		return &detailMethod{
			client:  issueClient,
			project: remote.RepositoryFullName(),
			id:      iid,
		}, nil
	}

	// Case of nothing Issue IID
	if createUpdateOption.hasEdit() {
		return &createOnEditorMethod{
			issueClient:      issueClient,
			repositoryClient: repositoryClient,
			opt:              createUpdateOption,
			project:          project,
			editFunc:         c.EditFunc,
		}, nil
	}
	if createUpdateOption.hasCreate() {
		return &createMethod{
			client:  issueClient,
			opt:     createUpdateOption,
			project: project,
		}, nil
	}
	if listOption.AllProject {
		return &listAllMethod{
			client: issueClient,
			opt:    listOption,
		}, nil
	}

	return &listMethod{
		client:  issueClient,
		opt:     listOption,
		project: project,
	}, nil
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
