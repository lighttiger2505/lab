package milestone

import (
	"bytes"
	"strings"

	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/internal/ui"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

const (
	ExitCodeOK    int = iota //0
	ExitCodeError int = iota //1
)

type Option struct {
	ProjectProfileOption *internal.ProjectProfileOption `group:"Project, Profile Options"`
	ListOption           *ListOption                    `group:"List Options"`
}

type ListOption struct {
	Num int `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of milestone to output."`
}

func newOptionParser(opt *Option) *flags.Parser {
	opt.ProjectProfileOption = &internal.ProjectProfileOption{}
	opt.ListOption = &ListOption{}
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = `config - Edit and Show config

Synopsis:
  # Edit config
  lab config

  # Show config
  lab config -l`
	return parser
}

type MilestoneCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	ClientFactory   lab.APIClientFactory
}

func (c *MilestoneCommand) Synopsis() string {
	return "List milestone"
}

func (c *MilestoneCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt Option
	parser := newOptionParser(&opt)
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *MilestoneCommand) Run(args []string) int {
	var opt Option
	parser := newOptionParser(&opt)
	_, err := parser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	pInfo, err := c.RemoteCollecter.CollectTarget(
		opt.ProjectProfileOption.Project,
		opt.ProjectProfileOption.Profile,
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if err := c.ClientFactory.Init(pInfo.ApiUrl(), pInfo.Token); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	client := c.ClientFactory.GetMilestoneClient()

	milestones, err := client.ListMilestones(
		pInfo.Project,
		makeListMilestoneOptions(opt.ListOption),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	output := milestoneOutput(milestones)
	result := columnize.SimpleFormat(output)
	c.UI.Message(result)

	return ExitCodeOK
}

func makeListMilestoneOptions(listOption *ListOption) *gitlab.ListMilestonesOptions {
	lopt := &gitlab.ListOptions{
		Page:    1,
		PerPage: listOption.Num,
	}
	opt := &gitlab.ListMilestonesOptions{
		ListOptions: *lopt,
	}
	return opt
}

func milestoneOutput(milestones []*gitlab.Milestone) []string {
	yellow := color.New(color.FgYellow).SprintFunc()
	var outputs []string
	for _, milestone := range milestones {
		output := strings.Join([]string{
			yellow(milestone.ID),
			milestone.Title,
			milestone.Description,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
