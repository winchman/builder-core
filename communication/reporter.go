package comm

import (
	"github.com/Sirupsen/logrus"
)

type Reporter struct {
	log   LogChan
	event EventChan
}

func NewReporter(log LogChan, event EventChan) *Reporter {
	return &Reporter{
		log:   log,
		event: event,
	}
}

// Log - send a log message into the ether
func (r *Reporter) Log(entry *logrus.Entry, message string) {
	entry.Message = message
	if r.log != nil {
		r.log <- NewLogEntry(entry)
	}
}

// LogLevel - send a log message into the ether, specifying level
func (r *Reporter) LogLevel(entry *logrus.Entry, message string, level logrus.Level) {
	entry.Level = level
	r.Log(entry, message)
}

// EventOptions are the options when telling a Reporter to trigger an event
type EventOptions struct {
	EventType EventType
	Note      string
}

// Event notifies the Reporter's EventChan that an event has occurred
func (r *Reporter) Event(opts EventOptions) {
	if r.event != nil {
		r.event <- &event{
			eventType: opts.EventType,
			note:      opts.Note,
		}
	}
}
