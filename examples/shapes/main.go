// Example shapes shows how to draw basic shapes into a window.
// It can be considered the Go aequivalent of
// https://x.org/releases/X11R7.5/doc/libxcb/tutorial/#drawingprim
// Four points, a single polyline, two line segments,
// two rectangle and two arcs are drawn.
package main

import (
	"fmt"
	"unicode/utf16"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

func main() {
	X, err := xgb.NewConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer X.Close()

	setup := xproto.Setup(X)
	screen := setup.DefaultScreen(X)
	wid, _ := xproto.NewWindowId(X)
	draw := xproto.Drawable(wid) // for now, we simply draw into the window

	// Create the window
	xproto.CreateWindow(X, screen.RootDepth, wid, screen.Root,
		0, 0, 180, 200, 8, // X, Y, width, height, *border width*
		xproto.WindowClassInputOutput, screen.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{screen.WhitePixel, xproto.EventMaskStructureNotify | xproto.EventMaskExposure})

	// Map the window on the screen
	xproto.MapWindow(X, wid)

	// Up to here everything is the same as in the `create-window` example.
	// We opened a connection, created and mapped the window.
	// Note how this time the border width is set to 8 instead of 0.
	//
	// But this time we'll be drawing some basic shapes:

	// First of all we need to create a context to draw with.
	// The graphics context combines all properties (e.g. color, line width, font, fill style, ...)
	// that should be used to draw something. All available properties
	//
	// These properties can be set by or'ing their keys (xproto.Gc*)
	// and adding the value to the end of the values array.
	// The order in which the values have to be given corresponds to the order that they defined
	// mentioned in `xproto`.
	//
	// Here we create a new graphics context
	// which only has the foreground (color) value set to black:
	foreground, _ := xproto.NewGcontextId(X)
	mask := uint32(xproto.GcForeground)
	values := []uint32{screen.BlackPixel}
	xproto.CreateGC(X, foreground, draw, mask, values)

	// It is possible to set the foreground value to something different.
	// In production, this should use xorg color maps instead for compatibility
	// but for demonstration setting the color directly also works.
	// For more information on color maps, see the xcb documentation:
	// https://x.org/releases/X11R7.5/doc/libxcb/tutorial/#usecolor
	red, _ := xproto.NewGcontextId(X)
	mask = uint32(xproto.GcForeground)
	values = []uint32{0xff0000}
	xproto.CreateGC(X, red, draw, mask, values)

	// We'll create another graphics context that draws thick lines:
	thick, _ := xproto.NewGcontextId(X)
	mask = uint32(xproto.GcLineWidth)
	values = []uint32{10}
	xproto.CreateGC(X, thick, draw, mask, values)

	// It is even possible to set multiple properties at once.
	// Only remember to put the values in the same order as they're
	// defined in `xproto`:
	// Foreground is defined first, so we also set it's value first.
	// LineWidth comes second.
	blue, _ := xproto.NewGcontextId(X)
	mask = uint32(xproto.GcForeground | xproto.GcLineWidth)
	values = []uint32{0x0000ff, 4}
	xproto.CreateGC(X, blue, draw, mask, values)

	// Properties of an already created gc can also be changed
	// if the original values aren't needed anymore.
	// In this case, we will change the line width
	// and cap (line corner) style of our foreground context,
	// to smooth out the polyline:
	mask = uint32(xproto.GcLineWidth | xproto.GcCapStyle)
	values = []uint32{3, xproto.CapStyleRound}
	xproto.ChangeGC(X, foreground, mask, values)

	// Writing text needs a bit more setup -- we first have
	// to open the required font.
	// For all available fonts, install and run xfontsel.
	font, _ := xproto.NewFontId(X)
	fontname := "-gnu-unifont-*-*-*-*-16-*-*-*-*-*-*-*"
	err = xproto.OpenFontChecked(X, font, uint16(len(fontname)), fontname).Check()
	if err != nil {
		fmt.Println("Failed opening the font:", err)
		return
	}

	// And create a context from it. We simply pass the font's ID to the GcFont property.
	textCtx, _ := xproto.NewGcontextId(X)
	mask = uint32(xproto.GcForeground | xproto.GcBackground | xproto.GcFont)
	values = []uint32{screen.BlackPixel, screen.WhitePixel, uint32(font)}
	xproto.CreateGC(X, textCtx, draw, mask, values)
	text := convertStringToChar2b("Hellö World!") // Unicode capable!

	// Close the font handle:
	xproto.CloseFont(X, font)

	// After all, writing text is way more comfortable using Xft - it supports TrueType,
	// and overall better configuration.

	points := []xproto.Point{
		{X: 10, Y: 10},
		{X: 20, Y: 10},
		{X: 30, Y: 10},
		{X: 40, Y: 10},
	}

	// A polyline is essientially a line with multiple points.
	// The first point is placed absolutely inside the window,
	// while every other point is placed relative to the one before it.
	polyline := []xproto.Point{
		{X: 50, Y: 10},
		{X: 5, Y: 20},   // move 5 to the right, 20 down
		{X: 25, Y: -20}, // move 25 to the right, 20 up - notice how this point is level again with the first point
		{X: 10, Y: 10},  // move 10 to the right, 10 down
	}

	segments := []xproto.Segment{
		{X1: 100, Y1: 10, X2: 140, Y2: 30},
		{X1: 110, Y1: 25, X2: 130, Y2: 60},
		{X1: 0, Y1: 160, X2: 90, Y2: 100},
	}

	// Rectangles have a start coordinate (upper left) and width and height.
	rectangles := []xproto.Rectangle{
		{X: 10, Y: 50, Width: 40, Height: 20},
		{X: 80, Y: 50, Width: 10, Height: 40},
	}

	// This rectangle we will use to demonstrate filling a shape.
	rectangles2 := []xproto.Rectangle{
		{X: 150, Y: 50, Width: 20, Height: 60},
	}

	// Arcs are defined by a top left position (notice where the third line goes to)
	// their width and height, a starting and end angle.
	// Angles are defined in units of 1/64 of a single degree,
	// so we have to multiply the degrees by 64 (or left shift them by 6).
	arcs := []xproto.Arc{
		{X: 10, Y: 100, Width: 60, Height: 40, Angle1: 0 << 6, Angle2: 90 << 6},
		{X: 90, Y: 100, Width: 55, Height: 40, Angle1: 20 << 6, Angle2: 270 << 6},
	}

	for {
		evt, err := X.WaitForEvent()
		switch evt.(type) {
		case xproto.ExposeEvent:
			// Draw the four points we specified earlier.
			// Notice how we use the `foreground` context to draw them in black.
			// Also notice how even though we changed the line width to 3,
			// these still only appear as a single pixel.
			// To draw points that are bigger than a single pixel,
			// one has to either fill rectangles, circles or polygons.
			xproto.PolyPoint(X, xproto.CoordModeOrigin, draw, foreground, points)

			// Draw the polyline. This time we specified `xproto.CoordModePrevious`,
			// which means that every point is placed relatively to the previous.
			// If we were to use `xproto.CoordModeOrigin` instead,
			// we could specify each point absolutely on the screen.
			// It is also possible to use `xproto.CoordModePrevious` for drawing *points*
			// which means that each point would be specified relative to the previous one,
			// just as we did with the polyline.
			xproto.PolyLine(X, xproto.CoordModePrevious, draw, foreground, polyline)

			// Draw two lines in red.
			xproto.PolySegment(X, draw, red, segments)

			// Draw two thick rectangles.
			// The line width only specifies the width of the outline.
			// Notice how the second rectangle gets completely filled
			// due to the line width.
			xproto.PolyRectangle(X, draw, thick, rectangles)

			// Draw the circular arcs in blue.
			xproto.PolyArc(X, draw, blue, arcs)

			// There's also a fill variant for all drawing commands:
			xproto.PolyFillRectangle(X, draw, red, rectangles2)

			// Draw the text. Xorg currently knows two ways of specifying text:
			//  a) the (extended) ASCII encoding using ImageText8(..., []byte)
			//  b) UTF16 encoding using ImageText16(..., []Char2b) -- Char2b is simply two  bytes
			// at the bottom of this file, there are two utility functions that help
			// convert a go string into an array of Char2b's.
			xproto.ImageText16(X, byte(len(text)), draw, textCtx, 10, 160, text)

		case xproto.DestroyNotifyEvent:
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Char2b is defined as
// 	Byte1 byte
// 	Byte2 byte
// and is used as a utf16 character.
// This function takes a string and converts each rune into a char2b.
func convertStringToChar2b(s string) []xproto.Char2b {
	var chars []xproto.Char2b
	var p []uint16

	for _, r := range []rune(s) {
		p = utf16.Encode([]rune{r})
		if len(p) == 1 {
			chars = append(chars, convertUint16ToChar2b(p[0]))
		} else {
			// If the utf16 representation is larger than 2 bytes
			// we can not use it and insert a blank instead:
			chars = append(chars, xproto.Char2b{Byte1: 0, Byte2: 32})
		}
	}

	return chars
}

// convertUint16ToChar2b converts a uint16 (which is basically two bytes)
// into a Char2b by using the higher 8 bits of u as Byte1
// and the lower 8 bits of u as Byte2.
func convertUint16ToChar2b(u uint16) xproto.Char2b {
	return xproto.Char2b{
		Byte1: byte((u & 0xff00) >> 8),
		Byte2: byte((u & 0x00ff)),
	}
}
