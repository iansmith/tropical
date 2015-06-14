package std

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical"
)

type StdInteractor struct {
	children []*StdInteractor
}

//all the functions in package std should return the interface types of package
//tropical, if you want to prevent "tinkering" in internal state
func NewStdInteractor() tropical.Interactor {
	return &StdInteractor{
		children: make([]*StdInteractor, 10 /*initial capacity*/),
	}
}

func (s *StdInteractor) Children() []tropical.Interactor {
	//this is annoying, problem of contravariant return types or some such
	result := make([]tropical.Interactor, len(s.children))
	for i, child := range s.children {
		result[i] = child
	}
	return result
}

//RootInteractor connects the HTML Canvas object to the interactor tree.  If
//you treat the StdInteractor as Interactor, you can't see the implementation.
//Note the canvas is not exposed because if you want to muck with that you
//should be using the types provided by the tropical package.
type RootInteractor struct {
	tropical.Interactor
	c tropical.Canvas
}

//
// NewRootInteractor is the base parent for an interactor tree.  The htmlId provided
// must match the id of exactly one canvas element in the page. This type does
// not implement DrawSelf because he is the root of the drawing pass.
func NewRootInteractor(htmlId string) tropical.Interactor {
	i := NewStdInteractor().(*StdInteractor)
	c := NewCanvas(htmlId)
	return &RootInteractor{
		Interactor: i,
		c:          c,
	}
}

// Draw is the root of a drawing pass. This implementation offsets its children
// by 5px in x and y and draws a rounded rectangle in the provided background
// color before drawing it's children.
func (r *RootInteractor) Draw(fillColor string) {
	insetAmount := 5
	radius := 9

	r.c.Save()
	r.c.SetFillColor(fillColor)
	js.Global.Call("roundRect", r.c.(*canvasImpl).Context(), insetAmount, insetAmount,
		r.c.(*canvasImpl).Width()-insetAmount, r.c.(*canvasImpl).Height()-insetAmount,
		radius, true, false)
	r.c.Restore()
}
