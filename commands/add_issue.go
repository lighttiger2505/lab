package commands

import (
	"bytes"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var createIssueFlags CreateIssueFlags
var createIssueParser = flags.NewParser(&createIssueFlags, flags.Default)

type CreateIssueFlags struct {
	Title       string `short:"t" long:"title" description:"The title of an issue"`
	Description string `short:"d" long:"description" description:"The description of an issue"`
	AssigneeID  int    `short:"a" long:"assignee_id" description:"The ID of a user to assign issue"`
	MilestoneID int    `short:"m" long:"milestone_id" description:"The ID of a milestone to assign issue"`

	Labels string `short:"l" long:"labels" description:"Comma-separated label names for an issue"`
}

type AddIssueCommand struct {
	Ui ui.Ui
}

func (c *AddIssueCommand) Synopsis() string {
	return "Add issue"
}

func (c *AddIssueCommand) Help() string {
	buf := &bytes.Buffer{}
	createIssueParser.Usage = "add-issue [options]"
	createIssueParser.WriteHelp(buf)
	return buf.String()
}

func (c *AddIssueCommand) Run(args []string) int {
	if _, err := createIssueParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var title string
	var description string

	if createIssueFlags.Title == "" || createIssueFlags.Description == "" {
		message := createIssueMessage(createIssueFlags.Title, createIssueFlags.Description)

		editor, err := git.NewEditor("ISSUE", "issue", message)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		title, description, err = editor.EditTitleAndDescription()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		if editor != nil {
			defer editor.DeleteFile()
		}
	} else {
		title = createIssueFlags.Title
		description = createIssueFlags.Description
	}

	if title == "" {
		return ExitCodeOK
	}

	conf, err := config.NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := gitlab.GitlabRemote(c.Ui, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := gitlab.GitlabClient(c.Ui, gitlabRemote, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	createIssueOptions := &gitlabc.CreateIssueOptions{
		Title:       gitlabc.String(title),
		Description: gitlabc.String(description),
		AssigneeID:  nil,
		MilestoneID: nil,
		Labels:      []string{},
	}

	issue, _, err := client.Issues.CreateIssue(
		gitlabRemote.RepositoryFullName(),
		createIssueOptions,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	c.Ui.Message(fmt.Sprintf("#%d", issue.IID))

	return ExitCodeOK
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
