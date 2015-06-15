package tropical

import (
	"fmt"
)

type PickList interface {
	Hits() []Interactor
	AddHit(Interactor) []Interactor
	Len() int
}

type EventType int

const (
	MouseMove EventType = iota
	MouseUp
	MouseDown
)

func (e EventType) String() string {
	switch e {
	case MouseMove:
		return "MouseMove"
	case MouseUp:
		return "MouseUp"
	case MouseDown:
		return "MouseDown"
	default:
		panic(fmt.Sprintf("unknown event type %d", int(e)))
	}
}

type Event interface {
	Type() EventType
	X() int
	Y() int
	Translate(byX, byY int)
}

type MousePolicy interface {
	Process(event Event, root RootInteractor)
}

type Mouser interface {
	MouseUp(Event)
	MouseMove(Event)
	MouseDown(Event)
}
