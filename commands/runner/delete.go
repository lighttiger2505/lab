package runner

import (
	lab "github.com/lighttiger2505/lab/gitlab"
)

type deleteMethod struct {
	runnerClient lab.Runner
	project      string
	id           int
}

func (m *deleteMethod) Process() (string, error) {
	err := m.runnerClient.RemoveRunner(m.id)
	if err != nil {
		return "", err
	}
	return "", nil
}
