package lab

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
)

type Ui interface {
	Ask(string) (string, error)
	Say(string)
	Message(string)
	Error(string)
	Machine(string, ...string)
}

type BasicUi struct {
	Reader      io.Reader
	Writer      io.Writer
	ErrorWriter io.Writer
	l           sync.Mutex
	interrupted bool
	scanner     *bufio.Scanner
}

func NewBasicUi() *BasicUi {
	return &BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stdout,
	}
}

func (rw *BasicUi) Ask(query string) (string, error) {
	rw.l.Lock()
	defer rw.l.Unlock()

	if rw.interrupted {
		return "", errors.New("interrupted")
	}

	if rw.scanner == nil {
		rw.scanner = bufio.NewScanner(rw.Reader)
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	log.Printf("ui: ask: %s", query)
	if query != "" {
		if _, err := fmt.Fprint(rw.Writer, query+" "); err != nil {
			return "", err
		}
	}

	result := make(chan string, 1)
	go func() {
		var line string
		if rw.scanner.Scan() {
			line = rw.scanner.Text()
		}
		if err := rw.scanner.Err(); err != nil {
			log.Printf("ui: scan err: %s", err)
			return
		}
		result <- line
	}()

	select {
	case line := <-result:
		return line, nil
	case <-sigCh:
		// Print a newline so that any further output starts properly
		// on a new line.
		fmt.Fprintln(rw.Writer)

		// Mark that we were interrupted so future Ask calls fail.
		rw.interrupted = true

		return "", errors.New("interrupted")
	}
}

func (rw *BasicUi) Say(message string) {
	rw.l.Lock()
	defer rw.l.Unlock()

	log.Printf("ui: %s", message)
	_, err := fmt.Fprint(rw.Writer, message+"\n")
	if err != nil {
		log.Printf("[ERR] Failed to write to UI: %s", err)
	}
}

func (rw *BasicUi) Message(message string) {
	rw.l.Lock()
	defer rw.l.Unlock()

	log.Printf("ui: %s", message)
	_, err := fmt.Fprint(rw.Writer, message+"\n")
	if err != nil {
		log.Printf("[ERR] Failed to write to UI: %s", err)
	}
}

func (rw *BasicUi) Error(message string) {
	rw.l.Lock()
	defer rw.l.Unlock()

	writer := rw.ErrorWriter
	if writer == nil {
		writer = rw.Writer
	}

	log.Printf("ui error: %s", message)
	_, err := fmt.Fprint(writer, message+"\n")
	if err != nil {
		log.Printf("[ERR] Failed to write to UI: %s", err)
	}
}

func (rw *BasicUi) Machine(t string, args ...string) {
	log.Printf("machine readable: %s %#v", t, args)
}
