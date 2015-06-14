package tropical

type Interactor interface {
	Children() []Interactor
}

type DrawsSelf interface {
	Draw(c Canvas)
}

type DrawsChildren interface {
	DrawChildren(c Canvas, children []Interactor)
}

//
// Canvas is a thin veneer over the Html5 Canvas object
// Note that the HTML-level properties of Canvas are not exposed here.
//
//http://www.w3.org/TR/2dcontext/
//http://www.w3schools.com/tags/ref_canvas.asp

type Canvas interface {
	SetFillColor(rgbish string)
	FillRectangle(x, y, w, h int) //shorthand for defining a path and then filling it
	Save()                        //save clipping rect
	Restore()                     //restore clipping rect
	BeginPath()                   //start a path
	Rectangle(x, y, w, h int)     //rectangular path
	Clip()                        //set clipping rect
	Fill()                        //fill the current path
}
