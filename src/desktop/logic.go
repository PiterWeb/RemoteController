package desktop

import (
	"github.com/PiterWeb/RemoteController/src/net"
)

func connectToHost(offerEncoded string) string {

	answerResponse := make(chan string)

	go net.InitAnswer(offerEncoded, answerResponse)

	// msg := "ID copied to clipboard"

	return <-answerResponse

	// clipboard := me.wnd.Hwnd().OpenClipboard()
	// defer clipboard.CloseClipboard()
	// clipboard.EmptyClipboard()
	// clipboard.WriteString(<-answerResponse)
	// me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)

}

var answerResponse = make(chan string)

func createHost() string {

	offer := make(chan string)

	go net.InitOffer(offer, answerResponse)

	// msg := "ID copied to clipboard"
	// clipboard := me.wnd.Hwnd().OpenClipboard()
	// defer clipboard.CloseClipboard()
	// clipboard.EmptyClipboard()
	// clipboard.WriteString(<-offer)
	// me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)

	return <-offer

}

func connectToClient(response string) {

	answerResponse <- response

}
