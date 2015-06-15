package std

import (
	"github.com/iansmith/tropical"
)

type eventImpl struct {
	t    tropical.EventType
	x, y int
}

func (e *eventImpl) Type() tropical.EventType {
	return e.t
}
func (e *eventImpl) X() int {
	return e.x
}
func (e *eventImpl) Y() int {
	return e.y
}

func newEventImpl(t tropical.EventType, x, y int) tropical.Event {
	return &eventImpl{t, x, y}
}
