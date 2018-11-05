package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Client interface {
	// Lint
	Lint(content string) (*gitlab.LintResult, error)
}

type LabClient struct {
	Client *gitlab.Client
}

func NewLabClient(client *gitlab.Client) *LabClient {
	return &LabClient{Client: client}
}

func (l *LabClient) Lint(content string) (*gitlab.LintResult, error) {
	lintResult, _, err := l.Client.Validate.Lint(content)
	if err != nil {
		return nil, fmt.Errorf("Failed lint. Error: %s", err.Error())
	}
	return lintResult, nil
}

type MockLabClient struct {
	Client
	// Lint
	MockLint func(content string) (*gitlab.LintResult, error)
}

func (m *MockLabClient) Lint(content string) (*gitlab.LintResult, error) {
	return m.MockLint(content)
}
