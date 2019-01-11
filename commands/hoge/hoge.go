package hoge

import (
	"bytes"
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
)

type ProjectProfileOption struct {
	Project string `long:"project" value-name:"<title>" description:"Project"`
	Profile string `long:"profile" value-name:"<title>" description:"Profile"`
}

type Option struct {
	ProjectProfileOption *ProjectProfileOption `group:"Project/Profile Options"`
}

func newOptionParser(opt *Option) *flags.Parser {
	opt.ProjectProfileOption = &ProjectProfileOption{}
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "new"
	return parser
}

type HogeCommand struct {
	UI        ui.Ui
	Collecter *gitutil.RemoteCollecter
}

func (c *HogeCommand) Synopsis() string {
	return "hogehoge"
}

var opt Option
var parser = newOptionParser(&opt)

func (c *HogeCommand) Help() string {
	buf := &bytes.Buffer{}
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *HogeCommand) Run(args []string) int {
	if _, err := parser.ParseArgs(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load config, %s", err)
	}
	remoteCollecter := gitutil.NewRemoteCollecter(c.UI, cfg, git.NewGitClient())

	project := opt.ProjectProfileOption.Project
	profile := opt.ProjectProfileOption.Profile

	res, err := remoteCollecter.CollectTarget(project, profile)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	fmt.Println("domain, ", res)
	return 0
}
