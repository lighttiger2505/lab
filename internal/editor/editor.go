package editor

import (
	"os"
	"os/exec"
)

func OpenEditor(args ...string) error {
	editorEnv := os.Getenv("EDITOR")
	if editorEnv == "" {
		editorEnv = "vim"
	}

	c := exec.Command(editorEnv, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
