package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type PipelineCommandOption struct {
	PipelineOption *PipelineOption     `group:"Pipeline Options"`
	OutputOption   *ListPipelineOption `group:"Output Options"`
}

func newPipelineCommandParser(opt *PipelineCommandOption) *flags.Parser {
	opt.PipelineOption = newPipelineOption()
	opt.OutputOption = newListPipelineOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = `pipeline [options]

Synopsis:
  # List pipeline
  lab pipeline 

  # Show pipeline
  lab pipeline <Pipeline IID>
`
	return parser
}

type PipelineOption struct {
	Scope   string `short:"c" long:"scope" description:"The scope of pipelines, one of: running, pending, finished, branches, tags"`
	States  string `short:"t" long:"states" description:" The status of pipelines, one of: running, pending, success, failed, canceled, skipped"`
	OrderBy string `short:"o" long:"orderby" default:"id" default-mask:"id" description:"Order pipelines by id, status, ref, or user_id"`
}

func newPipelineOption() *PipelineOption {
	pipeline := flags.NewNamedParser("lab", flags.Default)
	pipeline.AddGroup("Pipeline Options", "", &PipelineOption{})
	return &PipelineOption{}
}

type ListPipelineOption struct {
	Num  int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of pipeline to output."`
	Sort string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print pipeline ordered in \"asc\" or \"desc\" order."`
}

func newListPipelineOption() *ListPipelineOption {
	return &ListPipelineOption{}
}

type PipelineOperation int

const (
	ListPipeline PipelineOperation = iota
	ShowPipeline
)

type PipelineCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *PipelineCommand) Synopsis() string {
	return "List pipeline, List pipeline jobs"
}

func (c *PipelineCommand) Help() string {
	var pipelineCommandOption PipelineCommandOption
	pipelineCommandParser := newPipelineCommandParser(&pipelineCommandOption)
	buf := &bytes.Buffer{}
	pipelineCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *PipelineCommand) Run(args []string) int {
	// Parse flags
	var pipelineCommandOption PipelineCommandOption
	pipelineCommandParser := newPipelineCommandParser(&pipelineCommandOption)
	parseArgs, err := pipelineCommandParser.ParseArgs(args)
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

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	operation := pipelineOperation(pipelineCommandOption, parseArgs)
	switch operation {
	case ListPipeline:
		pipelines, err := client.ProjectPipelines(
			gitlabRemote.RepositoryFullName(),
			makePipelineOptions(pipelineCommandOption.PipelineOption, pipelineCommandOption.OutputOption),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}

		result := columnize.SimpleFormat(pipelineOutput(pipelines))
		c.UI.Message(result)
	case ShowPipeline:
		pid, err := strconv.Atoi(parseArgs[0])
		if err != nil {
			fmt.Errorf("Invalid pipeline iid. value: %s, error: %s", parseArgs[0], err)
		}
		jobs, err := client.ProjectPipelineJobs(
			gitlabRemote.RepositoryFullName(),
			makePiplineJobsOptions(),
			pid,
		)
		result := columnize.SimpleFormat(jobOutput(jobs))
		c.UI.Message(result)
	}

	return ExitCodeOK
}

func pipelineOperation(opt PipelineCommandOption, args []string) PipelineOperation {
	if len(args) > 0 {
		return ShowPipeline
	}
	return ListPipeline
}

func makePipelineOptions(pipelineOption *PipelineOption, outputOption *ListPipelineOption) *gitlab.ListProjectPipelinesOptions {
	var scope *string
	if pipelineOption.Scope != "" {
		scope = gitlab.String(pipelineOption.Scope)
	}
	var status *gitlab.BuildStateValue
	if pipelineOption.States != "" {
		v := gitlab.BuildStateValue(pipelineOption.States)
		status = &v
	}
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: outputOption.Num,
	}
	listPipelinesOptions := &gitlab.ListProjectPipelinesOptions{
		Scope:       scope,
		Status:      status,
		Ref:         gitlab.String(""),
		YamlErrors:  gitlab.Bool(false),
		Name:        gitlab.String(""),
		Username:    gitlab.String(""),
		OrderBy:     gitlab.OrderBy(gitlab.OrderByValue(pipelineOption.OrderBy)),
		Sort:        gitlab.String(outputOption.Sort),
		ListOptions: *listOption,
	}
	return listPipelinesOptions
}

func makePiplineJobsOptions() *gitlab.ListJobsOptions {
	return &gitlab.ListJobsOptions{}
}

func pipelineOutput(pipelines gitlab.PipelineList) []string {
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

func jobOutput(jobs []*gitlab.Job) []string {
	var outputs []string
	for _, job := range jobs {
		output := strings.Join([]string{
			strconv.Itoa(job.ID),
			job.Status,
			job.Stage,
			job.Name,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
