package commands

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type PipelineCommand struct {
	UI       ui.Ui
	Provider gitlab.Provider
}

func (c *PipelineCommand) Synopsis() string {
	return "Show pipeline"
}

func (c *PipelineCommand) Help() string {
	buf := &bytes.Buffer{}
	return buf.String()
}

func (c *PipelineCommand) Run(args []string) int {
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

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	pipelines, err := client.ProjectPipelines(gitlabRemote.RepositoryFullName())
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	result := columnize.SimpleFormat(pipelineOutput(pipelines))
	c.UI.Message(result)

	return ExitCodeOK
}

func pipelineOutput(pipelines gitlabc.PipelineList) []string {
	var outputs []string
	for _, pipeline := range pipelines {
		output := strings.Join([]string{
			strconv.Itoa(pipeline.ID),
			pipeline.Status,
			pipeline.Ref,
			pipeline.Sha,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
