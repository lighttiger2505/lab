package runner

import (
	"strconv"
	"strings"

	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type listMethod struct {
	runnerClient lab.Runner
	opt          *ListOption
	project      string
}

func (m *listMethod) Process() (string, error) {
	runners, err := m.runnerClient.ListProjectRunners(
		m.project,
		makeListProjectRunnerOptions(m.opt),
	)
	if err != nil {
		return "", err
	}
	result := columnize.SimpleFormat(listRunnerOutput(runners))
	return result, nil
}

func makeListRunnerOptions(opt *ListOption) *gitlab.ListRunnersOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: opt.Num,
	}
	listRunnersOptions := &gitlab.ListRunnersOptions{
		ListOptions: *listOption,
	}
	if opt.Scope != "" {
		listRunnersOptions.Scope = gitlab.String(opt.Scope)
	}
	return listRunnersOptions
}

func makeListProjectRunnerOptions(opt *ListOption) *gitlab.ListProjectRunnersOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: opt.Num,
	}
	listProjectRunnersOptions := &gitlab.ListProjectRunnersOptions{
		ListOptions: *listOption,
	}
	if opt.Scope != "" {
		listProjectRunnersOptions.Scope = gitlab.String(opt.Scope)
	}
	return listProjectRunnersOptions
}

func listRunnerOutput(runners []*gitlab.Runner) []string {
	var outputs []string
	for _, runner := range runners {
		output := strings.Join([]string{
			strconv.Itoa(runner.ID),
			runner.Name,
			runner.Description,
			runner.Status,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
