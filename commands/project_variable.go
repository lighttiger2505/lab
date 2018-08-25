package commands

import (
	"bytes"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type ProjectVaribleCommandOption struct {
	CreateUpdateOption *CreateUpdateProjectVaribleOption `group:"List Options"`
}

func newProjectVaribleOptionParser(opt *ProjectVaribleCommandOption) *flags.Parser {
	opt.CreateUpdateOption = newAddProjectVaribleOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = `project-variable - Create and Edit, list a project variable

Synopsis:
  # List project variables 
  lab project-variable

  # Create project variables 
  lab project-variable -a <key> <value>

  # Update project variables 
  lab project-variable -u <key> <value>
  # Remove project variables 

  # Show issue
  lab project-variable -d <key>`
	return parser
}

type ProjectVariableOperation int

const (
	CreateProjectVariable ProjectVariableOperation = iota
	UpdateProjectVariable
	RemoveProjectVariable
	ListProjectVariable
)

func projectVaribaleOperation(opt ProjectVaribleCommandOption, args []string) ProjectVariableOperation {
	createUpdateOption := opt.CreateUpdateOption

	if createUpdateOption.Add {
		return UpdateProjectVariable
	}
	if createUpdateOption.Update {
		return UpdateProjectVariable
	}
	if createUpdateOption.Delete {
		return RemoveProjectVariable
	}
	return ListProjectVariable
}

type CreateUpdateProjectVaribleOption struct {
	Add    bool `short:"a" long:"add" description:"Create/Add project variable."`
	Update bool `short:"u" long:"update" description:"Update project variable."`
	Delete bool `short:"d" long:"delete" description:"Delete project variable."`
}

func newAddProjectVaribleOption() *CreateUpdateProjectVaribleOption {
	return &CreateUpdateProjectVaribleOption{}
}

type ProjectVariableCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *ProjectVariableCommand) Synopsis() string {
	return "List project level variables"
}

func (c *ProjectVariableCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt ProjectVaribleCommandOption
	parser := newProjectVaribleOptionParser(&opt)
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *ProjectVariableCommand) Run(args []string) int {
	var opt ProjectVaribleCommandOption
	parser := newProjectVaribleOptionParser(&opt)
	parseArgs, err := parser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

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

	// Do issue operation
	switch projectVaribaleOperation(opt, parseArgs) {
	case ListProjectVariable:
		users, err := client.GetVariables(
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
		result := columnize.SimpleFormat(projectVariableOutput(users))
		c.UI.Message(result)
	}

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
