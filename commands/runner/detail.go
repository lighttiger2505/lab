package runner

import (
	"fmt"
	"strings"

	"github.com/lighttiger2505/lab/internal/api"
)

type detailMethod struct {
	runnerClient api.Runner
	id           int
}

func (m *detailMethod) Process() (string, error) {
	detail, err := m.runnerClient.GetRunnerDetails(m.id)
	if err != nil {
		return "", err
	}
	template := `%d
Status: %s
Description: %s
Token: %s
Tag: %s
Version :%s
AccessLevel: %s
MaximumTimeout: %d
`
	res := fmt.Sprintf(
		template,
		detail.ID,
		detail.Status,
		detail.Description,
		detail.Token,
		strings.Join(detail.TagList, ", "),
		detail.Version,
		detail.AccessLevel,
		detail.MaximumTimeout,
	)
	return res, nil
}
