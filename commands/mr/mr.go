package mr

import (
	"bytes"
	"fmt"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/browse"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type CreateUpdateOption struct {
	Edit         bool   `short:"e" long:"edit" description:"Edit the merge request on editor. Start the editor with the contents in the given title and message options."`
	Title        string `short:"i" long:"title" value-name:"<title>" description:"The title of an merge request"`
	Message      string `short:"m" long:"message" value-name:"<message>" description:"The message of an merge request"`
	Template     string `short:"p" long:"template" value-name:"<merge request template>" description:"The template of an merge request"`
	SourceBranch string `long:"source" value-name:"<source branch>" description:"The source branch"`
	TargetBranch string `long:"target" value-name:"<target branch>" default:"master" default-mask:"master" description:"The target branch"`
	StateEvent   string `long:"state-event" value-name:"<state>" description:"Change the status. \"opened\", \"closed\""`
	AssigneeID   int    `long:"assignee-id" value-name:"<assignee id>" description:"The ID of assignee."`
}

type ListOption struct {
	Num        int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of merge request to output."`
	State      string `long:"state" value-name:"<state>" default:"all" default-mask:"all" description:"Print only merge request of the state just those that are \"opened\", \"closed\", \"merged\" or \"all\""`
	Scope      string `long:"scope" value-name:"<scope>" default:"all" default-mask:"all" description:"Print only given scope. \"created-by-me\", \"assigned-to-me\" or \"all\"."`
	OrderBy    string `long:"orderby" value-name:"<orderby>" default:"updated_at" default-mask:"updated_at" description:"Print merge request ordered by \"created_at\" or \"updated_at\" fields."`
	Sort       string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print merge request ordered in \"asc\" or \"desc\" order."`
	Search     string `short:"s" long:"search"  value-name:"<search word>" description:"Search merge request against their title and description."`
	Milestone  string `long:"milestone"  value-name:"<milestone>" description:"lists merge request that have an assigned milestone."`
	AuthorID   int    `long:"author-id"  value-name:"<auther id>" description:"lists merge request that have an author id."`
	AssigneeID int    `long:"assignee-id"  value-name:"<assignee id>" description:"lists merge request that have an assignee id."`
	Opened     bool   `short:"o" long:"opened" description:"Shorthand of the state option for \"--state=opened\"."`
	Closed     bool   `short:"c" long:"closed" description:"Shorthand of the state option for \"--state=closed\"."`
	Merged     bool   `short:"g" long:"merged" description:"Shorthand of the state option for \"--state=merged\"."`
	CreatedMe  bool   `short:"r" long:"created-me" description:"Shorthand of the scope option for \"--scope=created-by-me\"."`
	AssignedMe bool   `short:"a" long:"assigned-me" description:"Shorthand of the scope option for \"--scope=assigned-by-me\"."`
	AllProject bool   `short:"A" long:"all-project" description:"Print the merge request of all projects"`
}

func (l *ListOption) getState() string {
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

func (l *ListOption) getScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

type ShowOption struct {
	NoComment bool `long:"no-comment" description:"Not print a list of comments for a spcific merge request."`
}

type BrowseOption struct {
	Browse bool `short:"b" long:"browse" description:"Browse merge request."`
}

type Option struct {
	ProjectProfileOption *internal.ProjectProfileOption `group:"Project, Profile Options"`
	CreateUpdateOption   *CreateUpdateOption            `group:"Create, Update Options"`
	ListOption           *ListOption                    `group:"List Options"`
	ShowOption           *ShowOption                    `group:"Show Options"`
	BrowseOption         *BrowseOption                  `group:"Browse Options"`
}

func newOptionParser(opt *Option) *flags.Parser {
	opt.ProjectProfileOption = &internal.ProjectProfileOption{}
	opt.CreateUpdateOption = &CreateUpdateOption{}
	opt.ListOption = &ListOption{}
	opt.BrowseOption = &BrowseOption{}
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = `merge-request - Create and Edit, list a merge request

Synopsis:
  # List merge request
  lab merge-request [-n <num>] -l [--state <state>] [--scope <scope>]
                    [--orderby <orderby>] [--sort <sort>] -o -c -g
                    -r -a -A

  # Create merge request
  lab merge-request [-e] [-i <title>] [-d <message>] [--assignee-id=<assignee id>]

  # Update merge request
  lab merge-request <merge request iid> [-t <title>] [-d <description>] [--state-event=<state>] [--assignee-id=<assignee id>]

  # Show merge request
  lab merge-request <mergerequest iid>
  
  # Browse merge request
  lab merge-request -b [<mergerequest iid>]`

	return parser
}

type MergeRequestCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	GitClient       git.Client
	ClientFactory   lab.APIClientFactory
	EditFunc        func(program, file string) error
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Create and Edit, list a merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt Option
	mergeRequestCommandParser := newOptionParser(&opt)
	mergeRequestCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	var opt Option
	mergeRequestCommandParser := newOptionParser(&opt)
	parseArgs, err := mergeRequestCommandParser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	pInfo, err := c.RemoteCollecter.CollectTarget(
		opt.ProjectProfileOption.Project,
		opt.ProjectProfileOption.Profile,
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if err := c.ClientFactory.Init(pInfo.ApiUrl(), pInfo.Token); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	method, err := c.getMethod(opt, parseArgs, pInfo, c.ClientFactory)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	res, err := method.Process()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if res != "" {
		c.UI.Message(res)
	}

	return ExitCodeOK
}

func (c *MergeRequestCommand) getMethod(opt Option, args []string, pInfo *gitutil.GitLabProjectInfo, clientFactory lab.APIClientFactory) (internal.Method, error) {
	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption
	browseOption := opt.BrowseOption
	showOption := opt.ShowOption

	mrClient := clientFactory.GetMergeRequestClient()
	repositoryClient := clientFactory.GetRepositoryClient()
	noteClient := clientFactory.GetNoteClient()

	iid, err := validMergeRequestIID(args)
	if err != nil {
		return nil, err
	}

	if browseOption.Browse {
		return &browseMethod{
			opener: &browse.Browser{},
			url:    pInfo.SubpageUrl("merge_requests"),
			id:     iid,
		}, nil
	}

	// Case of getting Merge Request IID
	if len(args) > 0 {
		if createUpdateOption.Edit {
			return &updateOnEditorMethod{
				client:   mrClient,
				opt:      createUpdateOption,
				project:  pInfo.Project,
				id:       iid,
				editFunc: c.EditFunc,
			}, nil
		}
		if hasCreateUpdateOption(createUpdateOption) {
			return &updateMethod{
				client:  mrClient,
				opt:     createUpdateOption,
				project: pInfo.Project,
				id:      iid,
			}, nil
		}

		return &detailMethod{
			mrClient:   mrClient,
			noteClient: noteClient,
			opt:        showOption,
			project:    pInfo.Project,
			id:         iid,
		}, nil
	}

	// Case of nothing MergeRequest IID
	if createUpdateOption.Edit {
		return &createOnEditorMethod{
			client:           mrClient,
			repositoryClient: repositoryClient,
			opt:              createUpdateOption,
			project:          pInfo.Project,
			editFunc:         c.EditFunc,
		}, nil

	}
	if createUpdateOption.Title != "" {
		return &createMethod{
			client:  mrClient,
			opt:     createUpdateOption,
			project: pInfo.Project,
		}, nil
	}

	if listOption.AllProject {
		return &listAllMethod{
			client: mrClient,
			opt:    listOption,
		}, nil

	}
	return &listMethod{
		client:  mrClient,
		opt:     listOption,
		project: pInfo.Project,
	}, nil
}

func hasCreateUpdateOption(opt *CreateUpdateOption) bool {
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
