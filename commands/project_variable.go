package commands

import (
	"strings"

	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type ProjectVariableCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *ProjectVariableCommand) Synopsis() string {
	return "List project level variables"
}

func (c *ProjectVariableCommand) Help() string {
	return ""
}

func (c *ProjectVariableCommand) Run(args []string) int {
	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	client, err := c.Provider.GetProjectVariableClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	users, err := client.GetVariables(
		gitlabRemote.RepositoryFullName(),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	result := columnize.SimpleFormat(projectVariableOutput(users))

	c.UI.Message(result)

	return ExitCodeOK
}
func projectVariableOutput(variables []*gitlab.ProjectVariable) []string {
	var outputs []string
	for _, variable := range variables {
		output := strings.Join([]string{
			variable.Key,
			variable.Value,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
