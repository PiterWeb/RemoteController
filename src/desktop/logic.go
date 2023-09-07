package desktop

import (
	"github.com/PiterWeb/RemoteController/src/net"
	"github.com/rodrigocfd/windigo/win/co"
)

var localCandidatesChan = make(chan string)
var remoteCandidatesChan = make(chan string)

func initLogic(me *MainWindow) {

	btnAddRemoteCandidates(me)
	btnCreateHostOnClick(me)
	btnClientConnectOnClick(me)
	copyLocalCandidates(me)

}

func btnClientConnectOnClick(me *MainWindow) {

	me.btnClient.On().BnClicked(func() {

		answerResponse := make(chan string)

		go net.InitAnswer(me.inputClient.Text(), answerResponse, localCandidatesChan, remoteCandidatesChan)

		msg := "ID copied to clipboard"

		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.WriteString(<-answerResponse)
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)

	})

}

func btnCreateHostOnClick(me *MainWindow) {

	answerResponse := make(chan string)

	me.createHost.On().BnClicked(func() {

		offer := make(chan string)

		go net.InitOffer(offer, answerResponse, localCandidatesChan, remoteCandidatesChan)

		msg := "ID copied to clipboard"
		clipboard := me.wnd.Hwnd().OpenClipboard()
		defer clipboard.CloseClipboard()
		clipboard.WriteString(<-offer)
		me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)

	})

	me.btnHost.On().BnClicked(func() {

		answerResponse <- me.inputHost.Text()

	})

}

func copyLocalCandidates(me *MainWindow) {

	// me.candidatesCopy.On().BnClicked(func() {

	// 	msg := "Candidates copied to clipboard"
	// 	clipboard := me.wnd.Hwnd().OpenClipboard()
	// 	defer clipboard.CloseClipboard()
	// 	clipboard.WriteString(<-localCandidatesChan)
	// 	me.wnd.Hwnd().MessageBox(msg, "Success", co.MB_ICONINFORMATION)


	// })

}

func btnAddRemoteCandidates(me *MainWindow) {

	// me.btnCandidates.On().BnClicked(func() {

	// 	remoteCandidatesChan <- me.inputCandidates.Text()

	// })

}
