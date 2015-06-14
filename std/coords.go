package std

import (
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

func NewCoords(x, y, w, h int) tropical.Coords {
	return &Coords{x, y, w, h}
}
