package tropical

//
// Coords gives a functional interface to a 2D coordinate system.
//
type Coords interface {
	X() int //relative to parent
	Y() int //relative to parent
	Width() int
	Height() int
	SetX(int)
	SetY(int)
	SetWidth(int)
	SetHeight(int)
}

//
// Manipulating the tree
//
type TreeManipulator interface {
	Children() []Interactor
	AppendChild(Interactor) []Interactor
	Parent() Interactor
}

//
// An interactor is a chinese menu of choices.
//
type Interactor interface {
	TreeManipulator
	Coords
}

//
// Root interactors know how to initiate traversals and know about the
// particulars of the display.
//
type RootInteractor interface {
	Interactor
	Pick(event Event) PickList
	Draw()
}

//
// XXX It's sucky from an API standpoint that the function DefaultDrawSelf in
// XXX std package isn't directly tied to implementing the DrawSelf API. This
// XXX means that in theory they could get out of sync.  Of course, in practice
// XXX thats going to blow up pretty quickly.  Same for DrawsChildren and
// XXX std.DefaultDrawChildren.

//DrawsSelf is the designation that the particular interactor knows how to draw
//itself.  The method Draw() is always called with the clipping rectangle and
//origin already set up properly (based on the child's requested x,y,w,h)
type DrawsSelf interface {
	DrawSelf(c Canvas)
}

//DrawChildren is the designation that a particular interactor knows how to draw
//its children.
type DrawsChildren interface {
	DrawChildren(c Canvas)
}

//PicksSelf indicates you want to do fancy pick handling.  Implementors are
//expected to keep the traversal going by continuing through their children.
type PicksSelf interface {
	PickSelf(Event, PickList) bool
}

//
// Canvas is a thin veneer over the Html5 Canvas object
// Note that the HTML-level properties of Canvas are not exposed here.
//
//http://www.w3.org/TR/2dcontext/
//http://www.w3schools.com/tags/ref_canvas.asp

type Canvas interface {
	SetFillColor(rgbish string)                            //set the fill color
	SetStrokeColor(rgbish string)                          //set the stroke color
	Save()                                                 //save clipping rect
	Restore()                                              //restore clipping rect
	BeginPath()                                            //start a path
	Rectangle(x, y, w, h int)                              //rectangular path
	Clip()                                                 //set clipping rect
	Fill()                                                 //fill the current path
	Stroke()                                               //stroke the current path
	MoveTo(x, y int)                                       //move to point
	LineTo(x, y int)                                       //from current point to this destination
	Translate(x, y int)                                    //change coord system
	DrawImageById(id string, x, y int)                     //you must insure that id is in the page
	Arc(x, y, radius int, startAngle, finishAngle float64) //do an arc, angle in radians from 3 o'clock, going counter-clockwise
	//composite functions
	FillRectangle(x, y, w, h int) //shorthand for defining a path and then filling it
	DrawLine(x1, y1, x2, y2 int)  //shorthand for defining a line path then stroking it
}
