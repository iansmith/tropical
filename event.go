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

//
// MouseDispatch dispatches mouse events to possibly multiple policies.
//
type MouseDispatch interface {
	Process(event Event, root RootInteractor)
}

//
// MousePolicy understands the particular semantics of a policy and how
// to focus that policy for a particular interactor.
type MousePolicy interface {
	Start(event Event, target Interactor) bool
	Process(event Event, target Interactor)
}

//
// MouseMonitor just allows you to keep track of where the mouse is, for things
// like special cursors and so forth.  The event is always in root coords.
//
type MouseMonitor interface {
	MouseUp(Event)
	MouseMove(Event)
	MouseDown(Event)
}

//
// CLICKER responds to a click inside and release inside of the boundary area
//
type Clicker interface {
	Click() //no event!
}

//
// DRAGGER responds to drag protocal.  The drag protocol sends the _offset_
// from start of drag as coords to Drag().  Drag() is called on start and
// and end as well, for making the visuals look nice().
//
type Dragger interface {
	DragStart()
	Drag(offsetX, offsetY int)
	DragEnd()
}
