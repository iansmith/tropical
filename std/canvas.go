package std

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"

	"github.com/iansmith/tropical"
)

//This implementation assumes that the browser is doing double buffering so there
//no need to do that ourselves.
//http://stackoverflow.com/questions/2795269/does-html5-canvas-support-double-buffering
type canvasImpl struct {
	element, context      *js.Object
	htmlWidth, htmlHeight int
}

func NewCanvas(elementName string) tropical.Canvas {
	elem := js.Global.Get("document").Call("getElementById", elementName)
	if elem == nil {
		panic(fmt.Sprint("your code and html are out of sync, missing reference: %s", elementName))
	}
	ctx := elem.Call("getContext", "2d")
	result := &canvasImpl{
		element:    elem,
		context:    ctx,
		htmlWidth:  elem.Get("width").Int(),
		htmlHeight: elem.Get("height").Int(),
	}
	return result
}

//
// DOM Level Methods
//
func (c *canvasImpl) Width() int {
	return c.htmlWidth
}
func (c *canvasImpl) Height() int {
	return c.htmlHeight
}

func (c *canvasImpl) Context() *js.Object {
	return c.context
}
func (c *canvasImpl) Element() *js.Object {
	return c.element
}

//
// Convenience Methods
//
func (c *canvasImpl) FillRectangle(x, y, w, h int) {
	c.BeginPath()
	c.Rectangle(x, y, w, h)
	c.Fill()
}
func (c *canvasImpl) DrawLine(x1, y1, x2, y2 int) {
	c.BeginPath()
	c.MoveTo(x1, y1)
	c.LineTo(x2, y2)
	c.Stroke()
}

//
// Pass through functions to the 2d drawing context
//
func (c *canvasImpl) Translate(x, y int) {
	c.context.Call("translate", x, y)
}
func (c *canvasImpl) MoveTo(x, y int) {
	c.context.Call("moveTo", x, y)
}
func (c *canvasImpl) LineTo(x, y int) {
	c.context.Call("lineTo", x, y)
}
func (c *canvasImpl) Save() {
	c.context.Call("save")
}
func (c *canvasImpl) Fill() {
	c.context.Call("fill")
}
func (c *canvasImpl) Stroke() {
	c.context.Call("stroke")
}
func (c *canvasImpl) Restore() {
	c.context.Call("restore")
}
func (c *canvasImpl) BeginPath() {
	c.context.Call("beginPath")
}
func (c *canvasImpl) Rectangle(x, y, w, h int) {
	c.context.Call("rect", x, y, w, h)
}
func (c *canvasImpl) Clip() {
	c.context.Call("clip")
}
func (c *canvasImpl) SetFillColor(rgbish string) {
	c.context.Call("setFillColor", rgbish)
}
func (c *canvasImpl) SetStrokeColor(rgbish string) {
	c.context.Set("strokeStyle", rgbish)
}
func (c *canvasImpl) Arc(x, y, radius int, startAngle, finishAngle float64) {
	c.context.Call("arc", x, y, radius, startAngle, finishAngle, false)
}
