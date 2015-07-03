package std

import (
	_ "fmt"

	"github.com/iansmith/tropical"
)

type DefaultMouseDispatch struct {
	focusPolicy       tropical.MousePolicy
	focusedInteractor tropical.Interactor

	FocusPolicies []tropical.MousePolicy
	Monitors      []tropical.MouseMonitor
}

func NewDefaultMouseDispatch() tropical.MouseDispatch {
	return &DefaultMouseDispatch{
		FocusPolicies: []tropical.MousePolicy{&ClickerPolicy{}, &DraggerPolicy{}},
		Monitors:      []tropical.MouseMonitor{},
	}
}

type ClickerPolicy struct {
}

type DraggerPolicy struct {
	startX, startY int
}

func (m *ClickerPolicy) Start(event tropical.Event, target tropical.Interactor) bool {
	_, ok := target.(tropical.Clicker)
	return ok && event.Type() == tropical.MouseDown
}

func (m *ClickerPolicy) Process(event tropical.Event, target tropical.Interactor) {
	focus := target.(tropical.Clicker) //this will panic if it's not a clickre!

	// we want to do a bounds check against the event so convert to local coords
	ToLocalFromGlobal(target, event)

	switch event.Type() {
	case tropical.MouseDown:
		//ignored, because we only get here if we have picked so there is nothing
		//do until the mouse up
	case tropical.MouseUp:
		if event.X() < 0 || event.Y() < 0 || event.X() >= target.Width() || event.Y() >= target.Height() {
			break //isn't inside the bounds
		}
		focus.Click()
	case tropical.MouseMove:
		//ignored
	default:
		print("clickerpolicy: unexpected event type ", event.Type().String(), "ignoring")
	}

	//restore event back to global coords
	ToGlobalFromLocal(target, event)
}

func (d *DraggerPolicy) Process(event tropical.Event, target tropical.Interactor) {
	focus := target.(tropical.Dragger) //this will panic if it's not a dragger!
	diffX := event.X() - d.startX
	diffY := event.Y() - d.startY
	switch event.Type() {
	case tropical.MouseDown:
		focus.DragStart()
		focus.Drag(diffX, diffY)
	case tropical.MouseUp:
		focus.Drag(diffX, diffY)
		focus.DragEnd()
	case tropical.MouseMove:
		focus.Drag(diffX, diffY)
	default:
		print("draggerpolicy: unexpected event type ", event.Type().String(), "ignoring")
	}

}

func (d *DraggerPolicy) Start(event tropical.Event, target tropical.Interactor) bool {
	_, ok := target.(tropical.Dragger)
	if ok && event.Type() == tropical.MouseDown {
		d.startX = event.X()
		d.startY = event.Y()
		return true
	}
	return false
}

func (d *DefaultMouseDispatch) Process(event tropical.Event, root tropical.RootInteractor) {
	//
	// focusedInteractor!=nil implies that somebody captured the start of a mouse
	// protocol and wanted to focus it on that object
	//
	if d.focusedInteractor != nil {
		//print(fmt.Sprintf("have focused object: %T wih policy %T", d.focusedInteractor, d.focusPolicy))
		if d.focusedInteractor != nil {
			d.focusPolicy.Process(event, d.focusedInteractor)
			if event.Type() == tropical.MouseUp {
				d.focusPolicy = nil
			}
		}
		if event.Type() == tropical.MouseUp {
			d.focusedInteractor = nil
		}
		return
	}

	list := root.Pick(event)
	/*
		for i, picked := range list.Hits() {
			print("pick list contains ", i, fmt.Sprintf("%T", picked))
		}
	*/

outer:
	for _, picked := range list.Hits() {
		for _, candidate := range d.FocusPolicies {
			if candidate.Start(event, picked) {
				d.focusedInteractor = picked
				d.focusPolicy = candidate
				d.focusPolicy.Process(event, picked)
				ToGlobalFromLocal(picked, event)
				//print(fmt.Sprintf("other root coords? %d,%d", event.X(), event.Y()))
				break outer
			}
		}
	}

	//allow monitors to also get the info
	for _, mon := range d.Monitors {
		switch event.Type() {
		case tropical.MouseDown:
			mon.MouseDown(event)
		case tropical.MouseUp:
			mon.MouseUp(event)
		case tropical.MouseMove:
			mon.MouseMove(event)
		default:
			print("monitor: unexpected event type ", event.Type().String(), "ignoring")
		}
	}
}
