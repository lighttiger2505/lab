package commands

import (
	"bytes"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type PipelineOpt struct {
	Line    int    `short:"n" long:"line" default:"20" default-mask:"20" description:"output the NUM lines"`
	Scope   string `short:"c" long:"scope" description:"The scope of pipelines, one of: running, pending, finished, branches, tags"`
	States  string `short:"t" long:"states" description:" The status of pipelines, one of: running, pending, success, failed, canceled, skipped"`
	OrderBy string `short:"o" long:"orderby" default:"id" default-mask:"id" description:"Order pipelines by id, status, ref, or user_id"`
	Sort    string `short:"s" long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order"`
}

var pipelineOptions PipelineOpt
var pipelineParser = flags.NewParser(&pipelineOptions, flags.Default)

type PipelineCommand struct {
	UI       ui.Ui
	Provider gitlab.Provider
}

func (c *PipelineCommand) Synopsis() string {
	return "Show pipeline"
}

func (c *PipelineCommand) Help() string {
	buf := &bytes.Buffer{}
	pipelineParser.Usage = "pipeline [options]"
	pipelineParser.WriteHelp(buf)
	return buf.String()
}

func (c *PipelineCommand) Run(args []string) int {
	// Parse flags
	if _, err := pipelineParser.ParseArgs(args); err != nil {
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

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	pipelines, err := client.ProjectPipelines(
		gitlabRemote.RepositoryFullName(),
		makePipelineOptions(pipelineOptions),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	result := columnize.SimpleFormat(pipelineOutput(pipelines))
	c.UI.Message(result)

	return ExitCodeOK
}

func makePipelineOptions(opt PipelineOpt) *gitlabc.ListProjectPipelinesOptions {
	var scope *string
	if opt.Scope != "" {
		scope = gitlabc.String(opt.Scope)
	}
	var status *gitlabc.BuildStateValue
	if opt.States != "" {
		v := gitlabc.BuildStateValue(opt.States)
		status = &v
	}
	listPipelinesOptions := &gitlabc.ListProjectPipelinesOptions{
		Scope:      scope,
		Status:     status,
		Ref:        gitlabc.String(""),
		YamlErrors: gitlabc.Bool(false),
		Name:       gitlabc.String(""),
		Username:   gitlabc.String(""),
		OrderBy:    gitlabc.OrderBy(gitlabc.OrderByValue(opt.OrderBy)),
		Sort:       gitlabc.String(opt.Sort),
	}
	return listPipelinesOptions
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
