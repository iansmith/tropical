package std

import (
	"github.com/iansmith/tropical"
)

type DefaultMousePolicy struct {
	focus        tropical.Mouser
	asInteractor tropical.Interactor
}

func NewDefaultMousePolicy() tropical.MousePolicy {
	return &DefaultMousePolicy{} //nil init for focus
}

func (d *DefaultMousePolicy) Process(event tropical.Event, root tropical.RootInteractor) {

	if d.focus != nil {
		ToLocalFromGlobal(d.asInteractor, event)
		switch event.Type() {
		case tropical.MouseUp:
			d.focus.MouseUp(event)
			d.focus = nil
			d.asInteractor = nil
		case tropical.MouseMove:
			d.focus.MouseMove(event)
		default:
			print("unexpected event type ", event.Type().String(), "ignoring")
		}
		return
	}

	//no focus, is this a mouse down?
	if event.Type() != tropical.MouseDown {
		print("not bothering with mouse event", event.Type().String())
		return
	}

	list := root.Pick(event)
	for _, picked := range list.Hits() {
		m, ok := picked.(tropical.Mouser)
		if !ok {
			continue
		}
		ToLocalFromGlobal(picked, event)
		m.MouseDown(event)
		d.focus = m
		d.asInteractor = picked
	}
}
