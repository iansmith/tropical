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

//
// Convenience Methods
//
func (c *canvasImpl) FillRectangle(x, y, w, h int) {
	c.Rectangle(x, y, w, h)
	c.Fill()
}

//
// Pass through functions to the 2d drawing context
//
func (c *canvasImpl) Save() {
	c.context.Call("save")
}
func (c *canvasImpl) Fill() {
	c.context.Call("fill")
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
