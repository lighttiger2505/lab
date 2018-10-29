package commands

import (
	"bytes"
	"fmt"
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
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
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
		return CreateProjectVariable
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

	op := projectVaribaleOperation(opt, parseArgs)
	if err := validProjectVariableArgs(op, parseArgs); err != nil {
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
	switch op {
	case CreateProjectVariable:
		_, err := client.CreateVariable(
			gitlabRemote.RepositoryFullName(),
			makeCreateProjectVariableOption(parseArgs[0], parseArgs[1]),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
	case UpdateProjectVariable:
		_, err := client.UpdateVariable(
			gitlabRemote.RepositoryFullName(),
			parseArgs[0],
			makeUpdateProjectVariableOption(parseArgs[0], parseArgs[1]),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
	case RemoveProjectVariable:
		err := client.RemoveVariable(
			gitlabRemote.RepositoryFullName(),
			parseArgs[0],
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
	case ListProjectVariable:
		variables, err := client.GetVariables(
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
		result := columnize.SimpleFormat(projectVariableOutput(variables))
		c.UI.Message(result)
	}

	return ExitCodeOK
}

func validProjectVariableArgs(op ProjectVariableOperation, args []string) error {
	switch op {
	case CreateProjectVariable:
		if len(args) < 2 {
			return fmt.Errorf("Usage: lab project-variable -a <key> <value>")
		}
	case UpdateProjectVariable:
		if len(args) < 2 {
			return fmt.Errorf("Usage: lab project-variable -u <key> <value>")
		}
	case RemoveProjectVariable:
		if len(args) < 1 {
			return fmt.Errorf("Usage: lab project-variable -d <key>")
		}
	}
	return nil
}

func makeCreateProjectVariableOption(key, value string) *gitlab.CreateVariableOptions {
	opt := &gitlab.CreateVariableOptions{
		Key:   gitlab.String(key),
		Value: gitlab.String(value),
	}
	return opt
}

func makeUpdateProjectVariableOption(key, value string) *gitlab.UpdateVariableOptions {
	opt := &gitlab.UpdateVariableOptions{
		Key:   gitlab.String(key),
		Value: gitlab.String(value),
	}
	return opt
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
