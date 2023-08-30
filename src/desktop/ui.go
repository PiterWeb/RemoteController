package desktop

import (
	"image/color"

	"gioui.org/io/system"
	"gioui.org/text"

	"gioui.org/layout"
	"gioui.org/op"

	"gioui.org/app"
	"gioui.org/widget/material"
)

func styleWindowUI(w *app.Window) {
	w.Option(app.Title("RemoteController"))
	w.Option(app.Size(500, 800))
}

func frameUI(e system.FrameEvent, ops *op.Ops, th *material.Theme) {
	setWindowBackground(th)
	gtx := layout.NewContext(ops, e)
	drawTitle(gtx, th)
	e.Frame(gtx.Ops)
}

func setWindowBackground(th *material.Theme) {

	*th = th.WithPalette(material.Palette{
		Bg: color.NRGBA{
			R: 8, G: 103, B: 126, A: 255,
		},
	})

}

func drawTitle(gtx layout.Context, th *material.Theme) {

	const titleText = "Remote Controller"

	title := material.H1(th, titleText)
	textColor := color.NRGBA{R: 121, G: 212, B: 253, A: 255}
	title.Color = textColor
	title.Alignment = text.Middle
	title.Layout(gtx)
}
