package pipeline

import (
	"bytes"
	"fmt"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type Option struct {
	ProjectProfileOption *internal.ProjectProfileOption `group:"Project, Profile Options"`
	ListOption           *ListOption                    `group:"List Options"`
	BrowseOption         *BrowseOption                  `group:"Brwose Options"`
}

func newOptionParser(opt *Option) *flags.Parser {
	opt.ProjectProfileOption = &internal.ProjectProfileOption{}
	opt.ListOption = &ListOption{}
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = `pipeline [options]

Synopsis:
  # List pipeline
  lab pipeline 

  # Show pipeline
  lab pipeline <Pipeline IID>
`
	return parser
}

type ListOption struct {
	Num     int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of pipeline to output."`
	Sort    string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print pipeline ordered in \"asc\" or \"desc\" order."`
	Scope   string `short:"c" long:"scope" description:"The scope of pipelines, one of: running, pending, finished, branches, tags"`
	States  string `short:"t" long:"states" description:" The status of pipelines, one of: running, pending, success, failed, canceled, skipped"`
	OrderBy string `short:"o" long:"orderby" default:"id" default-mask:"id" description:"Order pipelines by id, status, ref, or user_id"`
}

type BrowseOption struct {
	Browse bool `short:"b" long:"browse" description:"Browse issue."`
}

var opt Option
var parser = newOptionParser(&opt)

type PipelineCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	MethodFactory   MethodFactory
}

func (c *PipelineCommand) Synopsis() string {
	return "List pipeline, List pipeline jobs"
}

func (c *PipelineCommand) Help() string {
	buf := &bytes.Buffer{}
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *PipelineCommand) Run(args []string) int {
	parseArgs, err := parser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	iid, err := validIID(parseArgs)
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

	clientFacotry, err := lab.NewGitlabClientFactory(pInfo.ApiUrl(), pInfo.Token)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	method := c.MethodFactory.CreateMethod(opt, pInfo, iid, clientFacotry)
	res, err := method.Process()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if res != "" {
		c.UI.Message(res)
	}

	return ExitCodeOK
}

func validIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid args, please intput pipeline IID.")
	}
	return iid, nil
}
