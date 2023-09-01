package desktop

import (
	"github.com/rodrigocfd/windigo/ui"
	"github.com/rodrigocfd/windigo/win"
	"github.com/rodrigocfd/windigo/win/co"
)

// Creates a new instance of our main window.
func newWindow() *MyWindow {
	wnd := ui.NewWindowMain(
		ui.WindowMainOpts().
			Title("Remote Controller").
			ClientArea(win.SIZE{Cx: 500, Cy: 800}).WndStyles(co.WS_CAPTION | co.WS_SYSMENU | co.WS_CLIPCHILDREN | co.WS_VISIBLE | co.WS_MINIMIZEBOX | co.WS_MAXIMIZEBOX | co.WS_SIZEBOX).IconId(101),
	)

	me := &MyWindow{
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
		btnShow: ui.NewButton(wnd,
			ui.ButtonOpts().
				Text("&Connect").
				Position(win.POINT{X: 250, Y: 19}),
		),
	}

	me.btnShow.On().BnClicked(func() {
		const msg string = "Connection Stablished"
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)
	})

	return me
}
