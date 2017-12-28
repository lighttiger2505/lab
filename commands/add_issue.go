package commands

import (
	"fmt"
	"strings"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
)

type AddIssueCommand struct {
	Ui ui.Ui
}

func (c *AddIssueCommand) Synopsis() string {
	return "Add issue"
}

func (c *AddIssueCommand) Help() string {
	return "Usage: lab add-issue [option]"
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

	return ExitCodeOK
}
