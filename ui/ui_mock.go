package lab

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type MockUi struct {
	Reader      io.Reader
	Writer      *bytes.Buffer
	ErrorWriter *bytes.Buffer

	once sync.Once
}

func NewMockUi() *MockUi {
	return &MockUi{
		Writer:      new(bytes.Buffer),
		ErrorWriter: new(bytes.Buffer),
	}
}

func (u *MockUi) Ask(query string) (string, error) {
	var result string
	fmt.Fprint(u.Writer, query)
	if _, err := fmt.Fscanln(u.Reader, &result); err != nil {
		return "", err
	}

	return result, nil
}

func (u *MockUi) Say(message string) {
	fmt.Fprint(u.Writer, message+"\n")
}

func (u *MockUi) Message(message string) {
	fmt.Fprint(u.Writer, message+"\n")
}

func (u *MockUi) Error(message string) {
	fmt.Fprint(u.ErrorWriter, message+"\n")
}

func (u *MockUi) Machine(t string, args ...string) {
	return
}
