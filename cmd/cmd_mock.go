package cmd

import (
	"github.com/kballard/go-shellquote"
)

type MockCmd struct {
	Name string
	Args []string
}

func NewMockCmd(cmd string) *MockCmd {
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
	return &MockCmd{Name: name, Args: args}
}

func (cmd *MockCmd) SetCmd(name string) Cmd {
	cmd.Name = name
	return cmd
}

func (cmd *MockCmd) WithArg(arg string) Cmd {
	cmd.Args = append(cmd.Args, arg)
	return cmd
}

func (cmd *MockCmd) WithArgs(args ...string) Cmd {
	for _, arg := range args {
		cmd.WithArg(arg)
	}
	return cmd
}

func (cmd *MockCmd) CombinedOutput() (string, error) {
	return "", nil
}

func (cmd *MockCmd) Spawn() error {
	return nil
}
