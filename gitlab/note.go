package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Note interface {
	GetIssueNotes(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error)
}

type NoteClient struct {
	Note
	Client *gitlab.Client
}

func NewNoteClient(client *gitlab.Client) *NoteClient {
	return &NoteClient{Client: client}
}

func (c *NoteClient) GetIssueNotes(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error) {
	notes, _, err := c.Client.Notes.ListIssueNotes(repositoryName, iid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed get issue. %s", err.Error())
	}
	return notes, nil
}
