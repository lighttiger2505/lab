package cmd

import (
	"os"
	"os/exec"

	"github.com/kballard/go-shellquote"
)

type Cmd interface {
	SetCmd(name string) Cmd
	WithArg(arg string) Cmd
	WithArgs(args ...string) Cmd
	CombinedOutput() (string, error)
	Spawn() error
}

type BasicCmd struct {
	Name   string
	Args   []string
	Stdin  *os.File
	Stdout *os.File
	Stderr *os.File
}

func NewBasicCmd(cmd string) *BasicCmd {
	var name string
	var args []string

	cmds, _ := shellquote.Split(cmd)
	if len(cmds) > 0 {
		name = cmds[0]
		args = make([]string, 0)
		for _, arg := range cmds[1:] {
			args = append(args, arg)
		}
	}
	return &BasicCmd{Name: name, Args: args, Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
}

func (cmd *BasicCmd) SetCmd(name string) Cmd {
	cmd.Name = name
	return cmd
}

func (cmd *BasicCmd) WithArg(arg string) Cmd {
	cmd.Args = append(cmd.Args, arg)
	return cmd
}

func (cmd *BasicCmd) WithArgs(args ...string) Cmd {
	for _, arg := range args {
		cmd.WithArg(arg)
	}
	return cmd
}

func (cmd *BasicCmd) CombinedOutput() (string, error) {
	output, err := exec.Command(cmd.Name, cmd.Args...).CombinedOutput()
	return string(output), err
}

func (cmd *BasicCmd) Spawn() error {
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stdin = cmd.Stdin
	c.Stdout = cmd.Stdout
	c.Stderr = cmd.Stderr
	return c.Run()
}
