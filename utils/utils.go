package utils

import (
	"log"
	"os/exec"
	"strings"
)

func GitOutput(name string, args []string) string {
	return GitOutputs(name, args)[0]
}

func GitOutputs(name string, args []string) []string {
	var out = CmdOutput(name, args)
	var outs []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(line) != "" {
			outs = append(outs, string(line))
		}
	}
	return outs
}

func CmdOutput(name string, args []string) string {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
