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

		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.EmptyClipboard()
		clipboard.WriteString(<-answerResponse)
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)

	})

}

func btnCreateHostOnClick(me *MainWindow) {

	answerResponse := make(chan string)

	me.createClient.On().BnClicked(func() {

		offer := make(chan string)

		go net.InitOffer(offer, answerResponse)

		msg := "ID copied to clipboard"
		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.EmptyClipboard()
		clipboard.WriteString(<-offer)
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)

	})

	me.btnHost.On().BnClicked(func() {

		answerResponse <- me.inputHost.Text()

	})

}
