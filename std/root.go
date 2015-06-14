package std

import (
	_ "fmt"

	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical"
)

//RootInteractor connects the HTML Canvas object to the interactor tree.
//The canvas is not exposed because if you want to muck with that you
//should be using the types provided by the tropical package.  This type always
//has zero or one children.
type RootInteractor struct {
	tropical.Coords
	tropical.TreeManipulator
	c         tropical.Canvas
	fillColor string
}

//
// NewRootInteractor is the base parent for an interactor tree.  The htmlId provided
// must match the id of exactly one canvas element in the page. This type does
// not implement DrawSelf because he is the root of the drawing pass.  Child
// can be nil if you want to set it later.
func NewRootInteractor(htmlId string, fillColor string, child tropical.Interactor) tropical.Interactor {
	c := NewCanvas(htmlId).(*canvasImpl)
	result := &RootInteractor{
		c:               c,
		TreeManipulator: NewSingleChild(nil),
		Coords:          NewCoords(0, 0, c.Width(), c.Height()),
		fillColor:       fillColor,
	}
	return result
}

// Draw is the root of a drawing pass. This implementation offsets its children
// by 5px in x and y and draws a rounded rectangle in the provided background
// color before drawing its child.  If the child does not implement
// the DrawSelf() the default will be used.
func (r *RootInteractor) Draw() {
	insetAmount := 5
	radius := 9

	x := insetAmount
	y := insetAmount
	w := r.c.(*canvasImpl).Width() - insetAmount
	h := r.c.(*canvasImpl).Height() - insetAmount

	r.c.Save()
	r.c.SetFillColor(r.fillColor)
	js.Global.Call("roundRect", r.c.(*canvasImpl).Context(), x, y, w, h,
		radius, true, false)

	//setup the clipping rectangle to be the same path as our drawn rounded
	//rectangle so children can't be badly behaved
	r.c.BeginPath()
	js.Global.Call("roundRect", r.c.(*canvasImpl).Context(), x, y, w, h,
		radius, false, false)
	r.c.Clip()

	//setup our own coordinate system
	r.c.Translate(x, y)
	children := r.Children()
	if len(children) > 0 {
		d, ok := children[0].(tropical.DrawsSelf)
		if !ok {
			Default.DrawSelf(children[0], r.c)
		} else {
			d.DrawSelf(r.c)
		}
	}

	r.c.Restore()
}
