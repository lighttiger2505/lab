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

var createIssueCommandOption CreateIssueCommandOption
var createIssueCommandParser *flags.Parser = newCreateIssueCommandParser(&createIssueCommandOption)

type CreateIssueCommandOption struct {
	GlobalOpt *GlobalOption      `group:"Global Options"`
	CreateOpt *CreateIssueOption `group:"Create Options"`
}

func newCreateIssueCommandParser(opt *CreateIssueCommandOption) *flags.Parser {
	global := flags.NewNamedParser("lab", flags.Default)
	global.AddGroup("Global Options", "", &GlobalOption{})

	create := flags.NewNamedParser("lab", flags.Default)
	create.AddGroup("Create Issue Options", "", &CreateIssueOption{})

	opt.GlobalOpt = newGlobalOption()
	opt.CreateOpt = newCreateIssueOption()

	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "add-issue [options]"
	return parser
}

type CreateIssueOption struct {
	Title       string `short:"t" long:"title" description:"The title of an issue"`
	Description string `short:"d" long:"description" description:"The description of an issue"`
	AssigneeID  int    `short:"a" long:"assignee_id" description:"The ID of a user to assign issue"`
	MilestoneID int    `short:"m" long:"milestone_id" description:"The ID of a milestone to assign issue"`
	Labels      string `short:"l" long:"labels" description:"Comma-separated label names for an issue"`
}

func newCreateIssueOption() *CreateIssueOption {
	return &CreateIssueOption{}
}

type AddIssueCommand struct {
	Ui       ui.Ui
	Provider gitlab.Provider
}

func (c *AddIssueCommand) Synopsis() string {
	return "Add issue"
}

func (c *AddIssueCommand) Help() string {
	buf := &bytes.Buffer{}
	createIssueCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *AddIssueCommand) Run(args []string) int {
	if _, err := createIssueCommandParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	createOpt := createIssueCommandOption.CreateOpt
	title, description, err := getIssueTitleAndDesc(createOpt.Title, createOpt.Description)
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
	var gitlabRemote *git.RemoteInfo
	gOpt := issueOpt.GlobalOpt
	if gOpt.Repository != "" {
		namespace, project := gOpt.NameSpaceAndProject()
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

	issue, err := client.CreateIssue(
		makeCreateIssueOptions(title, description),
		gitlabRemote.RepositoryFullName(),
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	c.Ui.Message(fmt.Sprintf("#%d", issue.IID))

	return ExitCodeOK
}

func makeCreateIssueOptions(title, description string) *gitlabc.CreateIssueOptions {
	opt := &gitlabc.CreateIssueOptions{
		Title:       gitlabc.String(title),
		Description: gitlabc.String(description),
		MilestoneID: nil,
		Labels:      []string{},
	}
	return opt
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
