package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Repository interface {
	GetTree(repositoryName string, opt *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, error)
	GetFile(repositoryName string, filename string, opt *gitlab.GetRawFileOptions) (string, error)
}

type RepositoryClient struct {
	ProjectVariable
	Client *gitlab.Client
}

func NewRepositoryClient(client *gitlab.Client) Repository {
	return &RepositoryClient{Client: client}
}

func (c *RepositoryClient) GetTree(repositoryName string, opt *gitlab.ListTreeOptions) ([]*gitlab.TreeNode, error) {
	res, _, err := c.Client.Repositories.ListTree(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("failed list tree. %s", err.Error())
	}
	return res, nil
}

func (c *RepositoryClient) GetFile(repositoryName string, filename string, opt *gitlab.GetRawFileOptions) (string, error) {
	res, _, err := c.Client.RepositoryFiles.GetRawFile(repositoryName, filename, opt)
	if err != nil {
		return "", fmt.Errorf("failed get row file. %s", err.Error())
	}
	return string(res), nil
}
