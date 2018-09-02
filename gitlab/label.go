package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Label interface {
	GetLabels(repositoryName string, opt *gitlab.ListLabelsOptions) ([]*gitlab.Label, error)
}

type LabelClient struct {
	Label
	Client *gitlab.Client
}

func NewLabelClient(client *gitlab.Client) Label {
	return &LabelClient{Client: client}
}

func (c *LabelClient) GetLabels(repositoryName string, opt *gitlab.ListLabelsOptions) ([]*gitlab.Label, error) {
	res, _, err := c.Client.Labels.ListLabels(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("failed get row file. %s", err.Error())
	}
	return res, nil
}
