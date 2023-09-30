package desktop

import (
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

type MainWindow struct {
	wnd          ui.WindowMain
	createClient ui.Button
	inputHost    ui.Edit
	btnHost      ui.Button
	orTxt        ui.Static
	inputClient  ui.Edit
	btnClient    ui.Button
}

// Creates a new instance of our main window.
func initWindow() *MainWindow {
	wnd := ui.NewWindowMain(
		ui.WindowMainOpts().
			Title("Remote Controller").
			ClientArea(win.SIZE{Cx: 500, Cy: 800}).
			WndStyles(co.WS_CAPTION | co.WS_SYSMENU | co.WS_CLIPCHILDREN | co.WS_VISIBLE | co.WS_MINIMIZEBOX | co.WS_SIZEBOX).
			IconId(101),
	)

	me := &MainWindow{
		wnd: wnd,
		createClient: ui.NewButton(wnd,
			ui.ButtonOpts().
				Text("&Create Client").
				Position(win.POINT{X: 10, Y: 39}).
				Size(win.SIZE{Cx: 100}),
		),
		inputHost: ui.NewEdit(wnd,
			ui.EditOpts().
				Position(win.POINT{X: 190, Y: 40}).
				Size(win.SIZE{Cx: 150}),
		),
		btnHost: ui.NewButton(wnd,
			ui.ButtonOpts().
				Text("&Connect to Host").
				Position(win.POINT{X: 350, Y: 39}).
				Size(win.SIZE{Cx: 100}),
		),
		orTxt: ui.NewStatic(wnd,
			ui.StaticOpts().
				Text("OR").
				Position(win.POINT{X: 225, Y: 70}),
		),

		inputClient: ui.NewEdit(wnd,
			ui.EditOpts().
				Position(win.POINT{X: 90, Y: 96}).
				Size(win.SIZE{Cx: 150}),
		),
		btnClient: ui.NewButton(wnd,
			ui.ButtonOpts().
				Text("&Connect to Client").
				Position(win.POINT{X: 250, Y: 95}).
				Size(win.SIZE{Cx: 100}),
		),
	}

	return me
}
