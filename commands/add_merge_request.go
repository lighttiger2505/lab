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

var createMergeReqeustFlags CreateMergeReqeustFlags
var createMergeReqeustParser = flags.NewParser(&createMergeReqeustFlags, flags.Default)

type CreateMergeReqeustFlags struct {
	Title        string `short:"t" long:"title" description:"The title of an merge request"`
	Description  string `short:"d" long:"description" description:"The description of an merge request"`
	SourceBranch string `short:"s" long:"source" description:"The source branch"`
	TargetBranch string `short:"g" long:"target" description:"The target branch"`
	AssigneeID   int    `short:"a" long:"assignee_id" description:"The ID of a user to assign merge request"`
	MilestoneID  int    `short:"m" long:"milestone_id" description:"The ID of a milestone to assign merge request"`

	Labels string `short:"l" long:"labels" description:"Comma-separated label names for an merge request"`
}

type AddMergeReqeustCommand struct {
	Ui ui.Ui
}

func (c *AddMergeReqeustCommand) Synopsis() string {
	return "Add merge request"
}

func (c *AddMergeReqeustCommand) Help() string {
	buf := &bytes.Buffer{}
	createMergeReqeustParser.Usage = "add-merge-request [options]"
	createMergeReqeustParser.WriteHelp(buf)
	return buf.String()
}

func (c *AddMergeReqeustCommand) Run(args []string) int {
	if _, err := createMergeReqeustParser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Get merge request title and description
	// launch vim when non specific flags
	var title string
	var description string
	if createMergeReqeustFlags.Title == "" || createMergeReqeustFlags.Description == "" {
		cs := git.CommentChar()
		message := createMergeRequestMessage(createMergeReqeustFlags.Title, createMergeReqeustFlags.Description, cs)

		editor, err := git.NewEditor("MERGE_REQUEST", "merge-request", message)
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
		title = createMergeReqeustFlags.Title
		description = createMergeReqeustFlags.Description
	}

	if title == "" || description == "" {
		return ExitCodeOK
	}

	// Get source branch
	// current branch from local repository when non specific flags
	currentBranch, err := git.GitCurrentBranch()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	if createMergeReqeustFlags.SourceBranch != "" {
		currentBranch = createMergeReqeustFlags.SourceBranch
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

	createMergeReqeustOptions := &gitlabc.CreateMergeRequestOptions{
		Title:           gitlabc.String(title),
		Description:     gitlabc.String(description),
		SourceBranch:    gitlabc.String(currentBranch),
		TargetBranch:    gitlabc.String(createMergeReqeustFlags.TargetBranch),
		AssigneeID:      nil,
		TargetProjectID: nil,
	}

	mergeRequest, _, err := client.MergeRequests.CreateMergeRequest(
		gitlabRemote.RepositoryFullName(),
		createMergeReqeustOptions,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	c.Ui.Message(fmt.Sprintf("#%d", mergeRequest.IID))

	return ExitCodeOK
}

func createMergeRequestMessage(title, description, cs string) string {
	message := strings.Replace(`%s
# Creating an merge request

# Write a message for this merge request. The first block of
# text is the title and the rest is the description.
%s
`, "#", cs, -1)
	message = fmt.Sprintf(message, title, description)
	return message
}
