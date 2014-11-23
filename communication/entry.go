package comm

import (
	"strings"

	"github.com/Sirupsen/logrus"
)

//type State uint8

//const (
//Requested State = iota
//Building
//Pushing
//Completed
//Errored
//)

type (
	// LogChan is a channel for log entries
	LogChan chan LogEntry

	// StatusChan is a channel for status updates (somewhat TBD)
	StatusChan chan StatusEntry

	// ExitChan is a channel for receiving the final exit value (error or nil)
	ExitChan chan error
)

// LogEntry is a convenient, extensible way to ship logrus log entries
type LogEntry interface {
	Entry() *logrus.Entry
}

// NewLogEntry produces a log entry that can be sent on a LogChan
func NewLogEntry(l *logrus.Entry) LogEntry {
	return (*logEntry)(l)
}

type logEntry logrus.Entry

func (l *logEntry) Entry() *logrus.Entry {
	return (*logrus.Entry)(l)
}

// StatusEntry is the tentative interface for the structs returned on the status channel
type StatusEntry interface {
	BuildID() int
	//PreviousState() State
	//CurrentState() State
	Note() string
	//Error() error // should be checked for non-nil
}

// NewLogEntryWriter returns a log entry writer initialized with the provided
// channel. The provided writer implements the io.Writer interface
func NewLogEntryWriter(ch LogChan) *LogEntryWriter {
	return &LogEntryWriter{log: ch}
}

// LogEntryWriter is a type for implementing the io.Writer interface
type LogEntryWriter struct {
	log LogChan
}

func (writer *LogEntryWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		writer.log <- (*logEntry)(&logrus.Entry{Message: string(line)})
	}
	return len(p), nil
}
