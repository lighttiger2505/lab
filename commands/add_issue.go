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
	"strings"
)

var createIssueFlags CreateIssueFlags
var createIssueParser = flags.NewParser(&createIssueFlags, flags.Default)

type CreateIssueFlags struct {
	Title       int    `short:"t" long:"title" description:"The title of an issue"`
	Description string `short:"d" long:"description" description:"The description of an issue"`
	AssigneeID  string `short:"a" long:"assignee_id" description:"The ID of a user to assign issue"`
	MilestoneID string `short:"m" long:"milestone_id" description:"The ID of a milestone to assign issue"`
	Labels      string `short:"l" long:"labels" description:"Comma-separated label names for an issue"`
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
	var title string
	var body string

	cs := git.CommentChar()
	message := strings.Replace(`
# Creating an issue
#
# Write a message for this issue. The first block of
# text is the title and the rest is the description.
`, "#", cs, -1)

	editor, err := git.NewEditor("ISSUE", "issue", message)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	title, body, err = editor.EditTitleAndBody()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	c.Ui.Message(fmt.Sprintf("title=%s, body=%s", title, body))

	if editor != nil {
		defer editor.DeleteFile()
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
		Title:       &title,
		Description: &body,
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
