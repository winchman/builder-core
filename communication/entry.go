package comm

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

type State uint8

const (
	Requested State = iota
	Building
	Pushing
	Completed
	Errored
)

type (
	LogChan    chan LogEntry
	StatusChan chan StatusEntry
	ExitChan   chan error
)

type LogEntry interface {
	Entry() *logrus.Entry
}

func NewLogEntry(l *logrus.Entry) LogEntry {
	return (*logEntry)(l)
}

type logEntry logrus.Entry

func (l *logEntry) Entry() *logrus.Entry {
	return (*logrus.Entry)(l)
}

// StatusMessage is the tentative interface for the structs returned on the status channel
type StatusEntry interface {
	BuildID() int
	PreviousState() State
	CurrentState() State
	Note() string
	//Error() error // should be checked for non-nil
}

// NewLogEntryWriter returns a log entry writer initialized with the provided
// channel. The provided writer implements the io.Writer interface
func NewLogEntryWriter(ch LogChan) *logEntryWriter {
	return &logEntryWriter{log: ch}
}

type logEntryWriter struct {
	log LogChan
}

func (writer *logEntryWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		writer.log <- (*logEntry)(&logrus.Entry{Message: string(line)})
	}
	return len(p), nil
}
