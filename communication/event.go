package comm

type EventType uint8

type Event interface {
	EventType() EventType
	Note() string
}

const (
	RequestedEvent EventType = iota

	//Building
	//Pushing
	//Completed
)

func (t EventType) String() string {
	switch t {
	case RequestedEvent:
		return "RequestedEvent"
	}
	return ""
}

type event struct {
	eventType EventType
	note      string
}

func (e *event) EventType() EventType { return e.eventType }
func (e *event) Note() string         { return e.note }
