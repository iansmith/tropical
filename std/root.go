package std

import (
	"fmt"
	"strconv"
	"strings"

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
func NewRootInteractor(htmlId string, fillColor string, child tropical.Interactor) (tropical.RootInteractor, chan tropical.Event) {
	c := NewCanvas(htmlId).(*canvasImpl)
	result := &RootInteractor{
		c:               c,
		TreeManipulator: NewSingleChild(nil),
		Coords:          NewPageCoords(htmlId),
		fillColor:       fillColor,
		eventChan:       make(chan tropical.Event),
	}

	c.Element().Set("onmousemove", func(e *js.Object) {
		x, y := result.fromPageToMyCoords(e.Get("pageX").Int(), e.Get("pageY").Int())
		go result.ProcessMouseEvent(newEventImpl(tropical.MouseMove, x, y))
	})
	c.Element().Set("onmousedown", func(e *js.Object) {
		x, y := result.fromPageToMyCoords(e.Get("pageX").Int(), e.Get("pageY").Int())
		go result.ProcessMouseEvent(newEventImpl(tropical.MouseDown, x, y))
	})
	c.Element().Set("onmouseup", func(e *js.Object) {
		x, y := result.fromPageToMyCoords(e.Get("pageX").Int(), e.Get("pageY").Int())
		go result.ProcessMouseEvent(newEventImpl(tropical.MouseUp, x, y))
	})
	return result, result.eventChan
}

func (r *RootInteractor) ProcessMouseEvent(e tropical.Event) {
	r.eventChan <- e
}

//XXX This may have an off-by-one error because it produces myY==height
//XXX instead of height-1 (which it should, since it produces myY==0)
//XXX Is this browser specific?
func (r *RootInteractor) fromPageToMyCoords(pageX, pageY int) (int, int) {
	x := pageX - r.X()
	y := pageY - r.Y()
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
// Pick() is the root of a Pick pass.  This implementation knows about the
// the fact that there is a border of insetAmount around the root interactor's
// actual drawing area.  Doesn't account for rounded corrners.
//
func (r *RootInteractor) Pick(event tropical.Event) tropical.PickList {
	trueX := event.X() - insetAmount
	trueY := event.Y() - insetAmount
	trueWidth := r.Width() - insetAmount
	trueHeight := r.Height() - insetAmount
	inDrawingArea := true
	pl := NewPickList()

	if trueX < 0 || trueY < 0 || trueX >= trueWidth || trueY >= trueHeight {
		inDrawingArea = false
	}
	if !inDrawingArea {
		return pl
	}
	if len(r.Children()) == 0 {
		return pl
	}
	child := r.Children()[0]
	event.Translate(insetAmount, insetAmount)
	p, ok := child.(tropical.PicksSelf)
	if !ok {
		Default.PickSelf(child, event, pl)
	} else {
		p.PickSelf(event, pl)
	}
	return pl
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

func parseBorderOrPanic(s string) int {
	parts := strings.Split(s, " ")
	if len(parts) < 2 || !strings.HasSuffix(parts[0], "px") {
		panic("can't understand border:" + s)
	}
	bd := strings.TrimSuffix(parts[0], "px")
	border, err := strconv.ParseInt(bd, 10, 32)
	if err != nil {
		panic("can't understand pixel amount in border:" + s)
	}
	return int(border)
}
func (p *PageCoords) X() int {
	bbox := p.elem.Call("getBoundingClientRect")
	x := bbox.Get("left").Int() + js.Global.Get("window").Get("scrollX").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	left := style.Call("getPropertyValue", "padding-left").Int()
	bd := style.Call("getPropertyValue", "border-left").String()
	border := parseBorderOrPanic(bd)
	return x + left + border
}
func (p *PageCoords) Y() int {
	bbox := p.elem.Call("getBoundingClientRect")
	y := bbox.Get("top").Int() + js.Global.Get("window").Get("scrollY").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	top := style.Call("getPropertyValue", "padding-top").Int()
	bd := style.Call("getPropertyValue", "border-top").String()
	border := parseBorderOrPanic(bd)
	return y + top + border
}
func (p *PageCoords) Width() int {
	bbox := p.elem.Call("getBoundingClientRect")
	w := bbox.Get("width").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	left := style.Call("getPropertyValue", "padding-left").Int()
	right := style.Call("getPropertyValue", "padding-right").Int()
	bd := style.Call("getPropertyValue", "border-left").String()
	bdL := parseBorderOrPanic(bd)
	bd = style.Call("getPropertyValue", "border-right").String()
	bdR := parseBorderOrPanic(bd)
	return w - (left + right + bdL + bdR)
}
func (p *PageCoords) Height() int {
	bbox := p.elem.Call("getBoundingClientRect")
	h := bbox.Get("height").Int()
	style := js.Global.Get("window").Call("getComputedStyle", p.elem, nil)
	top := style.Call("getPropertyValue", "padding-top").Int()
	bot := style.Call("getPropertyValue", "padding-bottom").Int()
	bd := style.Call("getPropertyValue", "border-top").String()
	bdT := parseBorderOrPanic(bd)
	bd = style.Call("getPropertyValue", "border-bottom").String()
	bdB := parseBorderOrPanic(bd)
	return h - (top + bot + bdT + bdB)
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
