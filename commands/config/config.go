package config

import (
	"bytes"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/internal/config"
	"github.com/lighttiger2505/lab/internal/editor"
	"github.com/lighttiger2505/lab/ui"
)

const (
	ExitCodeOK    int = iota //0
	ExitCodeError int = iota //1
)

type Option struct {
	List bool `short:"l" long:"list" description:"Show config."`
}

func newOptionParser(opt *Option) *flags.Parser {
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = `config - Edit and Show config

Synopsis:
  # Edit config
  lab config

  # Show config
  lab config -l`
	return parser
}

type ConfigCommand struct {
	UI     ui.UI
	Config *config.Config
}

func (c *ConfigCommand) Synopsis() string {
	return "Edit config"
}

func (c *ConfigCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt Option
	parser := newOptionParser(&opt)
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *ConfigCommand) Run(args []string) int {
	var opt Option
	parser := newOptionParser(&opt)
	_, err := parser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if opt.List {
		cfgStr, err := c.Config.Read()
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
		c.UI.Message(cfgStr)
	} else {
		editor.OpenEditor(c.Config.Path())
	}

	return ExitCodeOK
}
