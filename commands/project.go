package commands

import (
	"bytes"

	"github.com/lighttiger2505/lab/ui"
)

type ProjectCommand struct {
	UI ui.Ui
}

func (c *ProjectCommand) Synopsis() string {
	return "Show project"
}

func (c *ProjectCommand) Help() string {
	buf := &bytes.Buffer{}
	return buf.String()
}

func (c *ProjectCommand) Run(args []string) int {
	c.UI.Message("project command")
	return ExitCodeOK
}
