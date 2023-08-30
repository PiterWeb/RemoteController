package desktop

import (
	"os"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func frameLogic(e system.FrameEvent, ops *op.Ops) {

}

func InitWindow() {

	go func() {

		w := app.NewWindow()

		styleWindowUI(w)

		if err := runWindow(w); err != nil {
			panic(err)
		}

		os.Exit(0)

	}()

	app.Main()

}

func runWindow(w *app.Window) error {

	th := material.NewTheme()
	setWindowBackground(th)
	var ops op.Ops
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			frameUI(e, &ops, th)
			frameLogic(e, &ops)
		}
	}

}
