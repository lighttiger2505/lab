package api

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Note interface {
	GetIssueNotes(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error)
	GetMergeRequestNotes(repositoryName string, iid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error)
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
		return nil, fmt.Errorf("Failed get issue notes. %s", err.Error())
	}
	return notes, nil
}

func (c *NoteClient) GetMergeRequestNotes(repositoryName string, iid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error) {
	notes, _, err := c.Client.Notes.ListMergeRequestNotes(repositoryName, iid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed get merge request notes. %s", err.Error())
	}
	return notes, nil
}

type MockNoteClient struct {
	Note
	MockGetIssueNotes        func(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error)
	MockGetMergeRequestNotes func(repositoryName string, iid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error)
}

func (m *MockNoteClient) GetIssueNotes(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error) {
	return m.MockGetIssueNotes(repositoryName, iid, opt)
}

func (m *MockNoteClient) GetMergeRequestNotes(repositoryName string, iid int, opt *gitlab.ListMergeRequestNotesOptions) ([]*gitlab.Note, error) {
	return m.MockGetMergeRequestNotes(repositoryName, iid, opt)
}
