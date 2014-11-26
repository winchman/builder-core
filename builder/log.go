package builder

import (
	"github.com/Sirupsen/logrus"
	"github.com/sylphon/builder-core/communication"
)

// Log - send a log message into the ether
func (bob *Builder) Log(entry *logrus.Entry, message string) {
	entry.Message = message
	if bob.log != nil {
		bob.log <- comm.NewLogEntry(entry)
	}
}

// LogLevel - send a log message into the ether, specifying level
func (bob *Builder) LogLevel(entry *logrus.Entry, message string, level logrus.Level) {
	entry.Level = level
	bob.Log(entry, message)
}
