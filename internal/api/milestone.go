package api

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Milestone interface {
	ListMilestones(project string, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error)
}

type MilestoneClient struct {
	Client *gitlab.Client
}

func NewMilestoneClient(client *gitlab.Client) *MilestoneClient {
	return &MilestoneClient{Client: client}
}

func (c *MilestoneClient) ListMilestones(project string, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
	milestones, _, err := c.Client.Milestones.ListMilestones(project, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list milestone, %s", err.Error())
	}
	return milestones, nil
}

type MockMilestoneClient struct {
	MockListMilestones func(project string, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error)
}

func (m *MockMilestoneClient) ListMilestones(project string, opt *gitlab.ListMilestonesOptions) ([]*gitlab.Milestone, error) {
	return m.MockListMilestones(project, opt)
}
