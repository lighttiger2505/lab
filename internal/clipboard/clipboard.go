package clipboard

import (
	"fmt"

	"github.com/atotto/clipboard"
)

type Clipboard interface {
	Write(str string) error
	Read() (string, error)
}

type ClipboardRW struct{}

func (rw *ClipboardRW) Write(str string) error {
	if err := clipboard.WriteAll(str); err != nil {
		return fmt.Errorf("Error write %s to clipboard, %s", str, err)
	}
	return nil
}

func (rw *ClipboardRW) Read() (string, error) {
	res, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("Error read from clipboard, %s", err)
	}
	return res, nil
}

type MockClipboardRW struct {
	str string
}

func (rw *MockClipboardRW) Write(str string) error {
	rw.str = str
	return nil
}

func (rw *MockClipboardRW) Read() (string, error) {
	return rw.str, nil
}
