package api

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Lint interface {
	Lint(content string) (*gitlab.LintResult, error)
}

type LintClient struct {
	Client *gitlab.Client
}

func NewLintClient(client *gitlab.Client) *LintClient {
	return &LintClient{Client: client}
}

func (c *LintClient) Lint(content string) (*gitlab.LintResult, error) {
	lintResult, _, err := c.Client.Validate.Lint(content)
	if err != nil {
		return nil, fmt.Errorf("Failed lint. Error: %s", err.Error())
	}
	return lintResult, nil
}

type MockLintClient struct {
	MockLint func(content string) (*gitlab.LintResult, error)
}

func (m *MockLintClient) Lint(content string) (*gitlab.LintResult, error) {
	return m.MockLint(content)
}
