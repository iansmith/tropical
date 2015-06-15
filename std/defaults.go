package std

import (
	"github.com/iansmith/tropical"
)

//Defaults embodies the user-configurable default behaviors.  DrawSelf, for
//example, is accessed as std.Default.DrawSelf.  These are *functions* not
//objects, so you have to pass self to them.  The defaults are _free floating_
//and not associated with a particular type because we don't want to put
//any more requirements on the implementor of an interactor type.
type Defaults struct {
	DrawSelf        func(self tropical.Interactor, c tropical.Canvas)
	DrawChildren    func(self tropical.Interactor, c tropical.Canvas)
	StartDimensions func(self tropical.Interactor) (int, int)
	PickSelf        func(self tropical.Interactor, x, y int) tropical.PickList
}

var Default = &Defaults{}

func init() {
	//because of a initialization cycle, you must do this in an init() function
	Default.DrawSelf = DefaultDrawSelf
	Default.DrawChildren = DefaultDrawChildren
	Default.StartDimensions = DefaultStartDimensions
}

//DefaultDrawSelf is the standard implementation of DrawSelf for interactors
//that do not provide it.  The default is to draw nothing for itself and just
//try to draw any childen.  If this is a leaf node, this does nothing.
func DefaultDrawSelf(self tropical.Interactor, c tropical.Canvas) {
	if len(self.Children()) > 0 {
		dc, ok := self.(tropical.DrawsChildren)
		if !ok {
			Default.DrawChildren(self, c)
		} else {
			dc.DrawChildren(c)
		}
	}
}

//DefaultDrawChildren is the standard implementation of DrawChildren for interactors
//that do not provide it.  This implementation walks the children
//and sets up the clipping rectangle and the translation as requested by the child.
//It will call either the child's DrawSelf or the default one.
func DefaultDrawChildren(self tropical.Interactor, c tropical.Canvas) {
	for _, child := range self.Children() {
		childX := child.X()
		childY := child.Y()
		childW := child.Width()
		childH := child.Height()
		c.Save()
		c.BeginPath()
		c.Rectangle(childX, childY, childW, childH)
		c.Clip()
		c.Translate(childX, childY)
		d, ok := child.(tropical.DrawsSelf)
		if !ok {
			Default.DrawSelf(child, c)
		} else {
			d.DrawSelf(c)
		}
		c.Restore()
	}
}

//DefaultStartDimensions makes an interactor default to a visible, but small, size.
func DefaultStartDimensions(self tropical.Interactor) (int, int) {
	return 10, 10 //the old 10x10 trick!
}
