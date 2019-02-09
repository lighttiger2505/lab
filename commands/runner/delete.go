package runner

import "github.com/lighttiger2505/lab/internal/api"

type deleteMethod struct {
	runnerClient api.Runner
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
