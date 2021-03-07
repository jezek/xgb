// Example create-window shows how to create a window, map it, resize it,
// and listen to structure and key events (i.e., when the window is resized
// by the window manager, or when key presses/releases are made when the
// window has focus). The events are printed to stdout.
package main

import (
	"fmt"

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

	// xproto.Setup retrieves the Setup information from the setup bytes
	// gathered during connection.
	setup := xproto.Setup(X)

	// This is the default screen with all its associated info.
	screen := setup.DefaultScreen(X)

	// Any time a new resource (i.e., a window, pixmap, graphics context, etc.)
	// is created, we need to generate a resource identifier.
	// If the resource is a window, then use xproto.NewWindowId. If it's for
	// a pixmap, then use xproto.NewPixmapId. And so on...
	wid, _ := xproto.NewWindowId(X)

	// CreateWindow takes a boatload of parameters.
	xproto.CreateWindow(X, screen.RootDepth, wid, screen.Root,
		0, 0, 500, 500, 0,
		xproto.WindowClassInputOutput, screen.RootVisual, 0, []uint32{})

	// This call to ChangeWindowAttributes could be factored out and
	// included with the above CreateWindow call, but it is left here for
	// instructive purposes. It tells X to send us events when the 'structure'
	// of the window is changed (i.e., when it is resized, mapped, unmapped,
	// etc.) and when a key press or a key release has been made when the
	// window has focus.
	// We also set the 'BackPixel' to white so that the window isn't butt ugly.
	xproto.ChangeWindowAttributes(X, wid,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{ // values must be in the order defined by the protocol
			0xffffffff,
			xproto.EventMaskStructureNotify |
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease,
		})

	// MapWindow makes the window we've created appear on the screen.
	// We demonstrated the use of a 'checked' request here.
	// A checked request is a fancy way of saying, "do error handling
	// synchronously." Namely, if there is a problem with the MapWindow request,
	// we'll get the error *here*. If we were to do a normal unchecked
	// request (like the above CreateWindow and ChangeWindowAttributes
	// requests), then we would only see the error arrive in the main event
	// loop.
	//
	// Typically, checked requests are useful when you need to make sure they
	// succeed. Since they are synchronous, they incur a round trip cost before
	// the program can continue, but this is only going to be noticeable if
	// you're issuing tons of requests in succession.
	//
	// Note that requests without replies are by default unchecked while
	// requests *with* replies are checked by default.
	err = xproto.MapWindowChecked(X, wid).Check()
	if err != nil {
		fmt.Printf("Checked Error for mapping window %d: %s\n", wid, err)
	} else {
		fmt.Printf("Map window %d successful!\n", wid)
	}

	// This is an example of an invalid MapWindow request and what an error
	// looks like.
	err = xproto.MapWindowChecked(X, 0).Check()
	if err != nil {
		fmt.Printf("Checked Error for mapping window 0x1: %s\n", err)
	} else { // neva
		fmt.Printf("Map window 0x1 successful!\n")
	}

	// Start the main event loop.
	for {
		// WaitForEvent either returns an event or an error and never both.
		// If both are nil, then something went wrong and the loop should be
		// halted.
		//
		// An error can only be seen here as a response to an unchecked
		// request.
		ev, xerr := X.WaitForEvent()
		if ev == nil && xerr == nil {
			fmt.Println("Both event and error are nil. Exiting...")
			return
		}

		if ev != nil {
			fmt.Printf("Event: %s\n", ev)
		}
		if xerr != nil {
			fmt.Printf("Error: %s\n", xerr)
		}

		// This is how accepting events work:
		// The application checks what event we got
		// (the event must be registered using either xproto.CreateWindow (see l.35) or
		//  xproto.ChangeWindowAttributes (see l.50)).
		// and reacts to it accordingly. All events are defined in the xproto subpackage.
		switch ev.(type) {
		case xproto.KeyPressEvent:
			// See https://pkg.go.dev/github.com/jezek/xgb@v0.0.0-20210121230032-cec22bda1ce1/xproto#KeyPressEvent
			// for documentation about a key press event.
			kpe := ev.(xproto.KeyPressEvent)
			fmt.Printf("Key pressed: %d", kpe.Detail)
			// The Detail value depends on the keyboard layout,
			// for QWERTY, q is #24.
			if kpe.Detail == 24 {
				return // exit on q
			}
		case xproto.DestroyNotifyEvent:
			// Depending on the user's desktop environment (especially
			// window manager), killing a window might close the
			// client's X connection (e. g. the default Ubuntu
			// desktop environment).
			//
			// If that's the case for your environment, closing this example's window
			// will also call the underlying Go program (because closing the X
			// connection gives a nil event and EOF error).
			//
			// Consider how a single application might have multiple windows
			// (e.g. an open popup or dialog, ...)
			//
			// With other DEs, the X connection will still stay open even after the
			// X window is closed. For these DEs (e.g. i3) we have to check whether
			// the WM sent us a DestroyNotifyEvent and close our program.
			//
			// For more information about closing windows while maintaining
			// the X connection see
			// https://github.com/jezek/xgbutil/blob/master/_examples/graceful-window-close/main.go
			return
		}
	}
}
