package std

import (
	"github.com/iansmith/tropical"
)

//
// EVENT
//
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
func (e *eventImpl) Translate(byX, byY int) {
	e.x -= byX
	e.y -= byY
}

func newEventImpl(t tropical.EventType, x, y int) tropical.Event {
	return &eventImpl{t, x, y}
}

//
// PICK LIST
//

type pickListImpl struct {
	hits []tropical.Interactor
}

func (p *pickListImpl) Len() int {
	return len(p.hits)
}

func NewPickList() tropical.PickList {
	return &pickListImpl{
		[]tropical.Interactor{},
	}
}

func (p *pickListImpl) Hits() []tropical.Interactor {
	return p.hits
}

func (p *pickListImpl) AddHit(i tropical.Interactor) []tropical.Interactor {
	p.hits = append(p.hits, i)
	return p.hits
}
