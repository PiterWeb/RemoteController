package net

// import "testing"

// func TestInitAnswer(t *testing.T) {

// 	answerResponse := make(chan string)
// 	answerResponse2 := make(chan string)
// 	offerChan := make(chan string)
// 	triggerEnd := make(chan struct{})

// 	go InitOffer(offerChan, answerResponse2, triggerEnd)

// 	go InitAnswer(<-offerChan, answerResponse, triggerEnd)
// 	answerResponse2 <- (<-answerResponse)
// 	triggerEnd <- struct{}{}
// }
