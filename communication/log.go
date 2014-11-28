package comm

import (
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
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

// LogEntryWriter is a type for implementing the io.Writer interface
type LogEntryWriter interface {
	io.Writer
}

// NewLogEntryWriter returns a log entry writer initialized with the provided
// channel. The provided writer implements the io.Writer interface
func NewLogEntryWriter(ch LogChan) LogEntryWriter {
	return logEntryWriter{log: ch}
}

type logEntryWriter struct {
	log LogChan
}

func (writer logEntryWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		writer.log <- (*logEntry)(&logrus.Entry{Message: string(line), Level: logrus.DebugLevel})
	}
	return len(p), nil
}
