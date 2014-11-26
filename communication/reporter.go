package comm

import (
	"github.com/Sirupsen/logrus"
)

// Reporter is type for sending messages on log and/or status channels
type Reporter struct {
	log   LogChan
	event EventChan
}

// NewReporter returns a reporter that is initialized with the provided channels
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
	Data      map[string]interface{}
}

// Event notifies the Reporter's EventChan that an event has occurred
func (r *Reporter) Event(opts EventOptions) {
	if r.event != nil {
		r.event <- &event{
			eventType: opts.EventType,
			data:      opts.Data,
		}
	}
}
