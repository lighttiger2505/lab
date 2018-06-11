package commands

import (
	"bytes"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type JobCommandOption struct {
	ListOption *ListJobOption `group:"List Options"`
}

func newJobOptionParser(opt *JobCommandOption) *flags.Parser {
	opt.ListOption = newListJobOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "job [options]"
	parser.Usage = `issue - list a job

Synopsis:
  # List job
  lab issue [-n <num>] [--search=<search word>] [-A]`
	return parser
}

type ListJobOption struct {
	Num int `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of search to output."`
	// Scope string `long:"scope" value-name:"<scope>" default:"all" default-mask:"all" description:"Print only given scope. created, pending, running, failed, success, canceled, skipped, manual"`
}

func newListJobOption() *ListJobOption {
	return &ListJobOption{}
}

type JobCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *JobCommand) Synopsis() string {
	return "List job"
}

func (c *JobCommand) Help() string {
	var jobCommnadOption JobCommandOption
	jobCommnadOptionParser := newJobOptionParser(&jobCommnadOption)
	buf := &bytes.Buffer{}
	jobCommnadOptionParser.WriteHelp(buf)
	return buf.String()
}

func (c *JobCommand) Run(args []string) int {
	// Parse flags
	var jobCommnadOption JobCommandOption
	jobCommnadOptionParser := newJobOptionParser(&jobCommnadOption)
	if _, err := jobCommnadOptionParser.ParseArgs(args); err != nil {
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

	client, err := c.Provider.GetJobClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	listOpt := jobCommnadOption.ListOption
	var result string
	jobs, err := client.GetProjectJobs(
		makeProjectJobsOption(listOpt),
		gitlabRemote.RepositoryFullName(),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	result = columnize.SimpleFormat(projectJobOutput(jobs))

	c.UI.Message(result)

	return ExitCodeOK
}

func makeProjectJobsOption(opt *ListJobOption) *gitlab.ListJobsOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: opt.Num,
	}
	listJobOption := &gitlab.ListJobsOptions{
		ListOptions: *listOption,
		// Scope:       gitlab.String(opt.Scope),
	}
	return listJobOption
}

func projectJobOutput(jobs []gitlab.Job) []string {
	var outputs []string
	for _, job := range jobs {
		output := strings.Join([]string{
			strconv.Itoa(job.ID),
			job.Name,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
