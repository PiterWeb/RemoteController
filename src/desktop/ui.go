package desktop

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func InitWindow() {

	go func () {

		w:= app.NewWindow()
		if err := runWindow(w); err != nil {
			panic(err)
		}

		os.Exit(0)

	}()

	app.Main()

}

func runWindow(w *app.Window) error {

	th := material.NewTheme()
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:

			gtx := layout.NewContext(&ops, e)
			title := material.H1(th, "Hello, Gio")
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			title.Color = maroon
			title.Alignment = text.Middle
			title.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}

}