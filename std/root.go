package std

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical"
)

//RootInteractor connects the HTML Canvas object to the interactor tree.
//The canvas is not exposed because if you want to muck with that you
//should be using the types provided by the tropical package.  This type always
//has zero or one children. This type is kinda broken in that it does not account
//for the rounded rectangle corners of its display when doing event processing.
//Math is hard.
type RootInteractor struct {
	tropical.Coords
	tropical.TreeManipulator
	c         tropical.Canvas
	fillColor string
	eventChan chan tropical.Event
}

//border around our object
var insetAmount = 5

//
// NewRootInteractor is the base parent for an interactor tree.  The htmlId provided
// must match the id of exactly one canvas element in the page. This type does
// not implement DrawSelf because he is the root of the drawing pass.  Child
// can be nil if you want to set it later.
func NewRootInteractor(htmlId string, fillColor string, child tropical.Interactor) (tropical.Interactor, chan tropical.Event) {
	c := NewCanvas(htmlId).(*canvasImpl)
	result := &RootInteractor{
		c:               c,
		TreeManipulator: NewSingleChild(nil),
		Coords:          NewPageCoords(htmlId),
		fillColor:       fillColor,
		eventChan:       make(chan tropical.Event),
	}

	c.Element().Set("onmousemove", func(e *js.Object) {
		//x, y := result.fromPageToMyCoords(e.Get("pageX").Int(), e.Get("pageY").Int())
		//result.ProcessMouseEvent(newEventImpl(tropical.MouseMove, x, y))
	})
	c.Element().Set("onmousedown", func(e *js.Object) {
		x, y := result.fromPageToMyCoords(e.Get("pageX").Int(), e.Get("pageY").Int())
		result.ProcessMouseEvent(newEventImpl(tropical.MouseDown, x, y))
	})
	c.Element().Set("onmouseup", func(e *js.Object) {
		x, y := result.fromPageToMyCoords(e.Get("pageX").Int(), e.Get("pageY").Int())
		result.ProcessMouseEvent(newEventImpl(tropical.MouseUp, x, y))
	})
	return result, result.eventChan
}

func (r *RootInteractor) ProcessMouseEvent(e tropical.Event) {
	print("event", e.X(), e.Y(), e.Type().String(), "but", r.X(), r.Y(), r.Width(), r.Height())
}

//XXX This may have an off-by-one error because it produces myY==height
//XXX instead of height-1 (which it should, since it produces myY==0)
//XXX Is this browser specific?
func (r *RootInteractor) fromPageToMyCoords(pageX, pageY int) (int, int) {
	x := pageX - r.X()
	y := pageY - r.Y()

	//var style = window.getComputedStyle(document.getElementById("Example"), null);
	//style.getPropertyValue("height");
	style := js.Global.Get("window").Call("getComputedStyle", r.c.(*canvasImpl).Element(), nil)
	left := style.Call("getPropertyValue", "padding-left").Int()
	print("left padding", left, x-left)
	return x, y
}

// Draw is the root of a drawing pass. This implementation offsets its children
// by 5px in x and y and draws a rounded rectangle in the provided background
// color before drawing its child.  If the child does not implement
// the DrawSelf() the default will be used.
func (r *RootInteractor) Draw() {
	radius := 9

	x := insetAmount
	y := insetAmount
	w := r.Width() - insetAmount
	h := r.Height() - insetAmount

	print("root draw", x, y)
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

//
// PageCoords
//
type PageCoords struct {
	elem *js.Object
}

func NewPageCoords(htmlId string) tropical.Coords {
	elem := js.Global.Get("document").Call("getElementById", htmlId)
	if elem == nil {
		panic(fmt.Sprint("your code and html are out of sync, missing reference: %s", htmlId))
	}
	return &PageCoords{elem}
}

func (p *PageCoords) X() int {
	bbox := p.elem.Call("getBoundingClientRect")
	x := bbox.Get("left").Int() + js.Global.Get("window").Get("scrollX").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	left := style.Call("getPropertyValue", "padding-left").Int()
	return x + left
}
func (p *PageCoords) Y() int {
	bbox := p.elem.Call("getBoundingClientRect")
	y := bbox.Get("top").Int() + js.Global.Get("window").Get("scrollY").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	top := style.Call("getPropertyValue", "padding-top").Int()
	return y + top
}
func (p *PageCoords) Width() int {
	bbox := p.elem.Call("getBoundingClientRect")
	w := bbox.Get("width").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	left := style.Call("getPropertyValue", "padding-left").Int()
	right := style.Call("getPropertyValue", "padding-right").Int()
	return w - (left + right)
}
func (p *PageCoords) Height() int {
	bbox := p.elem.Call("getBoundingClientRect")
	h := bbox.Get("height").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	top := style.Call("getPropertyValue", "padding-top").Int()
	bot := style.Call("getPropertyValue", "padding-bottom").Int()
	return h - (top + bot)
}

func (p *PageCoords) SetX(int) {
	panic("cant call SetX on something controlled by HTML")
}
func (p *PageCoords) SetY(int) {
	panic("cant call SetY on something controlled by HTML")
}
func (p *PageCoords) SetWidth(int) {
	panic("cant call SetWidth on something controlled by HTML")
}
func (p *PageCoords) SetHeight(int) {
	panic("cant call SetWidth on something controlled by HTML")
}
