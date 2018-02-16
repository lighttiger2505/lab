package commands

import (
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var globalOpt GlobalOpt
var globalParser = flags.NewParser(&globalOpt, flags.Default)

type GlobalOpt struct {
	Repository string `short:"p" long:"repository" description:"target specific repository"`
}

func (g *GlobalOpt) ValidRepository() (string, string, error) {
	value := g.Repository
	splited := strings.Split(value, "/")
	if value != "" && len(splited) != 2 {
		return "", "", fmt.Errorf("Invalid repository \"%s\". Assumed input style is \"Namespace/Project\".", value)
	}
	return splited[0], splited[1], nil
}
