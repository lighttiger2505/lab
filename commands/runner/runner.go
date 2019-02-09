package runner

import (
	"bytes"
	"fmt"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/internal/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type Option struct {
	ProjectProfileOption *internal.ProjectProfileOption `group:"Project, Profile Options"`
	ListOption           *ListOption                    `group:"List Options"`
	DeleteOption         *DeleteOption                  `group:"Delete Options"`
}

func newParser(opt *Option) *flags.Parser {
	opt.ProjectProfileOption = &internal.ProjectProfileOption{}
	opt.ListOption = newListRunnerOption()
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "project [options]"
	return parser
}

type ListOption struct {
	Num   int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of runner to output."`
	Scope string `long:"scope" value-name:"<scope>" description:"Print only given scope. \"active\", \"paused\", \"online\" \"offline\"."`
}

type DeleteOption struct {
	Delete bool `short:"D" long:"delete" description:"delete registed runner."`
}

func newListRunnerOption() *ListOption {
	return &ListOption{}
}

type RunnerCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	ClientFactory   api.APIClientFactory
}

func (c *RunnerCommand) Synopsis() string {
	return "List CI/CD Runner"
}

var opt Option
var parser = newParser(&opt)

func (c *RunnerCommand) Help() string {
	buf := &bytes.Buffer{}
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *RunnerCommand) Run(args []string) int {
	parseArgs, err := parser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	id, err := validID(parseArgs)
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

	method := c.createMethod(id, opt, pInfo)
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

func (c *RunnerCommand) createMethod(id int, opt Option, pInfo *gitutil.GitLabProjectInfo) internal.Method {
	if id > 0 {
		if opt.DeleteOption.Delete {
			return &deleteMethod{
				runnerClient: c.ClientFactory.GetRunnerClient(),
				project:      pInfo.Project,
				id:           id,
			}
		}
		return &detailMethod{
			runnerClient: c.ClientFactory.GetRunnerClient(),
			id:           id,
		}
	}

	return &listMethod{
		runnerClient: c.ClientFactory.GetRunnerClient(),
		opt:          opt.ListOption,
		project:      pInfo.Project,
	}
}

func validID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid args, please input runner id.")
	}
	return id, nil
}
