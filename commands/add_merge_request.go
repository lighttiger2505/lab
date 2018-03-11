package commands

import (
	"bytes"
	"fmt"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var createMergeRequestCommandOption CreateMergeRequestCommandOption
var createMergeRequestCommandParser *flags.Parser = newCreateMergeRequestCommandParser(&createMergeRequestCommandOption)

type CreateMergeRequestCommandOption struct {
	GlobalOpt *GlobalOption             `group:"Global Options"`
	CreateOpt *CreateMergeRequestOption `group:"Create Options"`
}

func newCreateMergeRequestCommandParser(opt *CreateMergeRequestCommandOption) *flags.Parser {
	global := flags.NewNamedParser("lab", flags.Default)
	global.AddGroup("Global Options", "", &GlobalOption{})

	create := flags.NewNamedParser("lab", flags.Default)
	create.AddGroup("Create MergeRequest Options", "", &CreateMergeRequestOption{})

	opt.GlobalOpt = newGlobalOption()
	opt.CreateOpt = newCreateMergeRequestOption()

	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "add-merge-request [options]"
	return parser
}

type CreateMergeRequestOption struct {
	Title        string `short:"t" long:"title" description:"The title of an merge request"`
	Description  string `short:"d" long:"description" description:"The description of an merge request"`
	SourceBranch string `short:"s" long:"source" description:"The source branch"`
	TargetBranch string `short:"g" long:"target" description:"The target branch"`
	AssigneeID   int    `short:"a" long:"assignee_id" description:"The ID of a user to assign merge request"`
	MilestoneID  int    `short:"m" long:"milestone_id" description:"The ID of a milestone to assign merge request"`
	Labels       string `short:"l" long:"labels" description:"Comma-separated label names for an merge request"`
}

func newCreateMergeRequestOption() *CreateMergeRequestOption {
	return &CreateMergeRequestOption{}
}

type AddMergeReqeustCommand struct {
	Ui       ui.Ui
	Provider gitlab.Provider
}

func (c *AddMergeReqeustCommand) Synopsis() string {
	return "Add merge request"
}

func (c *AddMergeReqeustCommand) Help() string {
	buf := &bytes.Buffer{}
	createMergeRequestCommandParser.Usage = "add-merge-request [options]"
	createMergeRequestCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *AddMergeReqeustCommand) Run(args []string) int {
	if _, err := createMergeRequestCommandParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Get merge request title and description
	// launch vim when non specific flags
	createOpt := createMergeRequestCommandOption.CreateOpt
	title, description, err := getIssueTitleAndDesc(createOpt.Title, createOpt.Description)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Get source branch
	// current branch from local repository when non specific flags
	currentBranch, err := git.GitCurrentBranch()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	if createOpt.SourceBranch != "" {
		currentBranch = createOpt.SourceBranch
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

	mergeRequest, err := client.CreateMergeRequest(
		makeCreateMergeRequestOptios(createOpt, title, description, currentBranch),
		gitlabRemote.RepositoryFullName(),
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	c.Ui.Message(fmt.Sprintf("#%d", mergeRequest.IID))

	return ExitCodeOK
}

func makeCreateMergeRequestOptios(opt *CreateMergeRequestOption, title, description, branch string) *gitlabc.CreateMergeRequestOptions {
	createMergeRequestOption := &gitlabc.CreateMergeRequestOptions{
		Title:           gitlabc.String(title),
		Description:     gitlabc.String(description),
		SourceBranch:    gitlabc.String(branch),
		TargetBranch:    gitlabc.String(opt.TargetBranch),
		AssigneeID:      nil,
		TargetProjectID: nil,
	}
	return createMergeRequestOption
}

func createMergeRequestMessage(title, description string) string {
	message := `<!-- Write a message for this merge request. The first block of text is the title -->
%s

<!-- the rest is the description.  -->
%s
`
	message = fmt.Sprintf(message, title, description)
	return message
}

func getMergeRequestTitleAndDesc(titleIn, descIn string) (string, string, error) {
	var title, description string
	if titleIn == "" || descIn == "" {
		message := createMergeRequestMessage(titleIn, descIn)

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
