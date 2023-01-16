package eventing

type EventType int

const (
	Created EventType = iota
	Updated
	Deleted
)

type Event struct {
	Type             EventType
	Resource         interface{}
	PreviousResource interface{}
}

func NewEvent(t EventType, r interface{}, p interface{}) Event {
	return Event{t, r, p}
}
