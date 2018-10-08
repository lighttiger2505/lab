package mr

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

type detailMethod struct {
	internal.Method
	client  lab.MergeRequest
	project string
	id      int
}

func (m *detailMethod) Process() (string, error) {
	// Do get merge request
	mergeRequest, err := m.client.GetMergeRequest(m.id, m.project)
	if err != nil {
		return "", err
	}
	return outMergeRequestDetail(mergeRequest), nil
}

func outMergeRequestDetail(mergeRequest *gitlab.MergeRequest) string {
	base := `!%d
Title: %s
Assignee: %s
Author: %s
CreatedAt: %s
UpdatedAt: %s

%s`
	detial := fmt.Sprintf(
		base,
		mergeRequest.IID,
		mergeRequest.Title,
		mergeRequest.Assignee.Name,
		mergeRequest.Author.Name,
		mergeRequest.CreatedAt.String(),
		mergeRequest.UpdatedAt.String(),
		mergeRequest.Description,
	)
	return detial
}
