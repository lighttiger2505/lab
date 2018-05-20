package git

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func helperCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "git":
		switch args[0] {
		case "remote":
			remoteArgs := args[1:]
			if len(remoteArgs) < 1 {
				fmt.Println("origin\ngithub\ngitlab")
			}
			switch remoteArgs[0] {
			case "get-url":
				fmt.Println("git@gitlab.com:lighttiger2505/lab.git")
			default:
				fmt.Fprintf(os.Stderr, "Unknown remote args %v\n", args)
				os.Exit(2)
			}
		case "var":
			fmt.Println("vim")
		case "rev-parse":
			fmt.Println("/Users/lighttiger2505/dev/src/github.com/lighttiger2505/lab/.git")
		default:
			fmt.Fprintf(os.Stderr, "Unknown git command %v\n", args)
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(2)
	}
}

func TestRemoteInfos(t *testing.T) {
	execCommand = helperCommand
	defer func() { execCommand = exec.Command }()

	client := NewGitClient()
	results, err := client.RemoteInfos()
	if err != nil {
		t.Errorf("echo: %v", err)
	}

	got := results[0]
	want := &RemoteInfo{
		Remote:     "origin",
		Domain:     "gitlab.com",
		NameSpace:  "lighttiger2505",
		Repository: "lab",
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Invalid return value. want %v, got %v", want, got)
	}
}

func TestGitEditor(t *testing.T) {
	execCommand = helperCommand
	defer func() { execCommand = exec.Command }()

	got, err := GitEditor()
	if err != nil {
		t.Errorf("echo: %v", err)
	}

	want := "vim"
	if want != got {
		t.Errorf("Invalid return value. want %q, got %q", want, got)
	}
}

func TestGitDir(t *testing.T) {
	execCommand = helperCommand
	defer func() { execCommand = exec.Command }()

	got, err := GitDir()
	if err != nil {
		t.Errorf("echo: %v", err)
	}

	want := "/Users/lighttiger2505/dev/src/github.com/lighttiger2505/lab/.git"
	if want != got {
		t.Errorf("Invalid return value. want %q, got %q", want, got)
	}
}
