package mr

import (
	"bytes"
	"fmt"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

const MergeRequestTemplateDir = ".gitlab/merge_request_templates"

type CreateUpdateMergeRequestOption struct {
	Edit         bool   `short:"e" long:"edit" description:"Edit the merge request on editor. Start the editor with the contents in the given title and message options."`
	Title        string `short:"i" long:"title" value-name:"<title>" description:"The title of an merge request"`
	Message      string `short:"m" long:"message" value-name:"<message>" description:"The message of an merge request"`
	Template     string `short:"p" long:"template" value-name:"<merge request template>" description:"The template of an merge request"`
	SourceBranch string `short:"s" long:"source" description:"The source branch"`
	TargetBranch string `short:"t" long:"target" default:"master" default-mask:"master" description:"The target branch"`
	StateEvent   string `long:"state-event" description:"Change the status. \"opened\", \"closed\""`
	AssigneeID   int    `long:"assignee-id" description:"The ID of assignee."`
}

func newCreateUpdateMergeRequestOption() *CreateUpdateMergeRequestOption {
	return &CreateUpdateMergeRequestOption{}
}

type ListMergeRequestOption struct {
	Num        int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of merge request to output."`
	State      string `long:"state" value-name:"<state>" default:"all" default-mask:"all" description:"Print only merge request of the state just those that are \"opened\", \"closed\", \"merged\" or \"all\""`
	Scope      string `long:"scope" value-name:"<scope>" default:"all" default-mask:"all" description:"Print only given scope. \"created-by-me\", \"assigned-to-me\" or \"all\"."`
	OrderBy    string `long:"orderby" value-name:"<orderby>" default:"updated_at" default-mask:"updated_at" description:"Print merge request ordered by \"created_at\" or \"updated_at\" fields."`
	Sort       string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print merge request ordered in \"asc\" or \"desc\" order."`
	Opened     bool   `short:"o" long:"opened" description:"Shorthand of the state option for \"--state=opened\"."`
	Closed     bool   `short:"c" long:"closed" description:"Shorthand of the state option for \"--state=closed\"."`
	Merged     bool   `short:"g" long:"merged" description:"Shorthand of the state option for \"--state=merged\"."`
	CreatedMe  bool   `short:"r" long:"created-me" description:"Shorthand of the scope option for \"--scope=created-by-me\"."`
	AssignedMe bool   `short:"a" long:"assigned-me" description:"Shorthand of the scope option for \"--scope=assigned-by-me\"."`
	AllProject bool   `short:"A" long:"all-project" description:"Print the merge request of all projects"`
}

func (l *ListMergeRequestOption) GetState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	if l.Merged {
		return "merged"
	}
	return l.State
}

func (l *ListMergeRequestOption) GetScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

type MergeRequestOperation int

const (
	CreateMergeRequest MergeRequestOperation = iota
	CreateMergeRequestOnEditor
	UpdateMergeRequest
	UpdateMergeRequestOnEditor
	ShowMergeRequest
	ListMergeRequest
	ListMergeRequestAllProject
)

func newListMergeRequestOption() *ListMergeRequestOption {
	return &ListMergeRequestOption{}
}

type MergeRequestCommandOption struct {
	CreateUpdateOption *CreateUpdateMergeRequestOption `group:"Create, Update Options"`
	ListOption         *ListMergeRequestOption         `group:"List Options"`
}

func newMergeRequestOptionParser(opt *MergeRequestCommandOption) *flags.Parser {
	opt.CreateUpdateOption = newCreateUpdateMergeRequestOption()
	opt.ListOption = newListMergeRequestOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = `merge-request - Create and Edit, list a merge request

Synopsis:
  # List merge request
  lab merge-request [-n <num>] -l [--state <state>] [--scope <scope>]
                    [--orderby <orderby>] [--sort <sort>] -o -c -g
                    -r -a -A

  # Create merge request
  lab merge-request [-e] [-i <title>] [-d <message>] [--assignee-id=<assignee id>]

  # Update merge request
  lab merge-request <MergeRequest IID> [-t <title>] [-d <description>] [--state-event=<state>] [--assignee-id=<assignee id>]

  # Show merge request
  lab merge-request <MergeRequest IID>`
	return parser
}

type MergeRequestCommand struct {
	Ui        ui.Ui
	Provider  lab.Provider
	GitClient git.Client
	EditFunc  func(program, file string) error
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Create and Edit, list a merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	var mergeRequestCommandOption MergeRequestCommandOption
	mergeRequestCommandParser := newMergeRequestOptionParser(&mergeRequestCommandOption)
	mergeRequestCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	var mergeRequestCommandOption MergeRequestCommandOption
	mergeRequestCommandParser := newMergeRequestOptionParser(&mergeRequestCommandOption)
	parseArgs, err := mergeRequestCommandParser.ParseArgs(args)
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

	client, err := c.Provider.GetMergeRequestClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	iid, err := validMergeRequestIID(parseArgs)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	switch mergeRequestOperation(mergeRequestCommandOption, parseArgs) {
	case UpdateMergeRequest:
		method := &updateMethod{
			client:  client,
			opt:     mergeRequestCommandOption.CreateUpdateOption,
			project: gitlabRemote.RepositoryFullName(),
			id:      iid,
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		if res != "" {
			c.Ui.Message(res)
		}

	case UpdateMergeRequestOnEditor:
		method := &updateOnEditorMethod{
			client:   client,
			opt:      mergeRequestCommandOption.CreateUpdateOption,
			project:  gitlabRemote.RepositoryFullName(),
			id:       iid,
			editFunc: c.EditFunc,
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		if res != "" {
			c.Ui.Message(res)
		}

	case CreateMergeRequest:
		method := &createMethod{
			client:  client,
			opt:     mergeRequestCommandOption.CreateUpdateOption,
			project: gitlabRemote.RepositoryFullName(),
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case CreateMergeRequestOnEditor:
		repositoryClient, err := c.Provider.GetRepositoryClient(gitlabRemote)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		method := &createOnEditorMethod{
			client:           client,
			repositoryClient: repositoryClient,
			opt:              mergeRequestCommandOption.CreateUpdateOption,
			project:          gitlabRemote.RepositoryFullName(),
			editFunc:         c.EditFunc,
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case ShowMergeRequest:
		method := &detailMethod{
			client:  client,
			project: gitlabRemote.RepositoryFullName(),
			id:      iid,
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case ListMergeRequest:
		method := &listMethod{
			client:  client,
			opt:     mergeRequestCommandOption.ListOption,
			project: gitlabRemote.RepositoryFullName(),
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	case ListMergeRequestAllProject:
		method := &listAllMethod{
			client: client,
			opt:    mergeRequestCommandOption.ListOption,
		}
		res, err := method.Process()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(res)

	default:
		c.Ui.Error("Invalid merge request operation")
		return ExitCodeError
	}

	return ExitCodeOK
}

func mergeRequestOperation(opt MergeRequestCommandOption, args []string) MergeRequestOperation {
	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption

	// Case of getting Merge Request IID
	if len(args) > 0 {
		if createUpdateOption.Edit {
			return UpdateMergeRequestOnEditor
		}
		if hasCreateUpdateOption(createUpdateOption) {
			return UpdateMergeRequest
		}
		return ShowMergeRequest
	}

	// Case of nothing MergeRequest IID
	if createUpdateOption.Edit {
		return CreateMergeRequestOnEditor
	}
	if createUpdateOption.Title != "" {
		return CreateMergeRequest
	}
	if listOption.AllProject {
		return ListMergeRequestAllProject
	}

	return ListMergeRequest
}

func hasCreateUpdateOption(opt *CreateUpdateMergeRequestOption) bool {
	if opt.Title != "" || opt.Message != "" || opt.StateEvent != "" || opt.AssigneeID != 0 {
		return true
	}
	return false
}

func validMergeRequestIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid Issue IID. IID: %s", args[0])
	}
	return iid, nil
}

func editMergeRequestTemplate(title, description string) string {
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
