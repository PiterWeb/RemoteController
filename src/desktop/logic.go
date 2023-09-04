package desktop

import (
	"github.com/PiterWeb/RemoteController/src/net"
	"github.com/rodrigocfd/windigo/win/co"
)

func initLogic(me *MainWindow) {

	btnCreateHostOnClick(me)
	btnClientConnectOnClick(me)

}

func btnClientConnectOnClick(me *MainWindow) {

	me.btnClient.On().BnClicked(func() {

		answerResponse := make(chan string)

		go net.InitAnswer(me.inputClient.Text(), answerResponse)

		msg := "ID copied to clipboard"

		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)
		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.WriteString(<-answerResponse)

	})

}

func btnCreateHostOnClick(me *MainWindow) {

	answerResponse := make(chan string)

	me.createHost.On().BnClicked(func() {

		offer := make(chan string)

		go net.InitOffer(offer, answerResponse)

		msg := "ID copied to clipboard"
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)
		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.WriteString(<-offer)

	})

	me.btnHost.On().BnClicked(func() {

		answerResponse <- me.inputHost.Text()

	})

}
