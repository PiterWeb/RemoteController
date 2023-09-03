package desktop

import (
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

type MainWindow struct {
	wnd       ui.WindowMain
	lblName   ui.Static
	txtName   ui.Edit
	btnAnswer ui.Button
	btnOffer  ui.Button
}

// Creates a new instance of our main window.
func initWindow() *MainWindow {
	wnd := ui.NewWindowMain(
		ui.WindowMainOpts().
			Title("Remote Controller").
			ClientArea(win.SIZE{Cx: 500, Cy: 800}).
			WndStyles(co.WS_CAPTION | co.WS_SYSMENU | co.WS_CLIPCHILDREN | co.WS_VISIBLE | co.WS_MINIMIZEBOX | co.WS_MAXIMIZEBOX | co.WS_SIZEBOX).
			IconId(101),
	)

	me := &MainWindow{
		wnd: wnd,

		lblName: ui.NewStatic(wnd,
			ui.StaticOpts().
				Text("Connection ID").
				Position(win.POINT{X: 10, Y: 22}),
		),
		txtName: ui.NewEdit(wnd,
			ui.EditOpts().
				Position(win.POINT{X: 90, Y: 20}).
				Size(win.SIZE{Cx: 150}),
		),
		btnAnswer: ui.NewButton(wnd,
			ui.ButtonOpts().
				Text("&Connect").
				Position(win.POINT{X: 250, Y: 19}),
		),
		btnOffer: ui.NewButton(wnd,
			ui.ButtonOpts().
				Text("&Create ID").
				Position(win.POINT{X: 10, Y: 60}),
		),
	}

	return me
}
