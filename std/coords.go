package std

import (
	_ "fmt"

	"github.com/iansmith/tropical"
)

type Coords struct {
	x, y, w, h int
}

func (s *Coords) X() int {
	return s.x
}
func (s *Coords) Y() int {
	return s.y
}
func (s *Coords) Width() int {
	return s.w
}
func (s *Coords) Height() int {
	return s.h
}
func (s *Coords) SetX(x int) {
	s.x = x
}
func (s *Coords) SetY(y int) {
	s.y = y
}
func (s *Coords) SetWidth(width int) {
	s.w = width
}
func (s *Coords) SetHeight(height int) {
	s.h = height
}

func (s *Coords) ToLocalFromParent(parent tropical.Interactor, event tropical.Event) {
	event.Translate(s.X(), s.Y())
}
func (s *Coords) ToParentFromLocal(parent tropical.Interactor, event tropical.Event) {
	event.Translate(-s.X(), -s.Y())
}

func NewCoords(x, y, w, h int) tropical.Coords {
	return &Coords{x, y, w, h}
}

//if you pass a root interactor in here, nothing happens to the event (it's
//already global. otherwise we adjust the event to be in local coords of wanted.
func ToLocalFromGlobal(wanted tropical.Interactor, event tropical.Event) {
	parentChain := []tropical.Interactor{}
	curr := wanted
	for {
		if curr == nil {
			break
		}
		parentChain = append(parentChain, curr)
		curr = curr.Parent()
	}
	//don't need to adjust for the root object, since its already in root coords
	for i := 0; i < len(parentChain)-1; i++ {
		event.Translate(parentChain[i].X(), parentChain[i].Y())
	}
}
