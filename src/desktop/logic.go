package desktop

import (
	"fmt"

	"github.com/PiterWeb/RemoteController/src/net"
	"github.com/rodrigocfd/windigo/win/co"
)

func initLogic(me *MainWindow) {
	btnAnswerOnClick(me)
	btnOfferOnClick(me)
}

func btnAnswerOnClick(me *MainWindow) {

	me.btnAnswer.On().BnClicked(func() {

		go net.InitAnswer(me.txtName.Text())

		msg := fmt.Sprintf("Connection Stablished to ID %s", me.txtName.Text())
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)
	})

}

func btnOfferOnClick(me * MainWindow) {
	 
	me.btnOffer.On().BnClicked(func() {

		offer := make(chan string)

		go net.InitOffer(offer)

		msg := "ID copied to clipboard"
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)
		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.WriteString(<- offer)

	})

}
