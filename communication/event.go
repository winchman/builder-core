package comm

type EventType uint8

type Event interface {
	EventType() EventType
	Note() string
	Error() error // should be checked for non-nil
}

const (
	RequestedEvent EventType = iota

	//Building
	//Pushing
	//Completed
	ErrorEvent
)

type event struct {
	eventType EventType
	note      string
	err       error
}

func (e *event) EventType() EventType { return e.eventType }
func (e *event) Note() string         { return e.note }
func (e *event) Error() error         { return e.err }
