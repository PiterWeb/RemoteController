package net

import "testing"

func TestInitAnswer(t *testing.T) {

	answerResponse := make(chan string)
	answerResponse2 := make(chan string)
	offerChan := make(chan string)

	go InitOffer(offerChan, answerResponse2)

	go InitAnswer(<-offerChan, answerResponse)
	answerResponse2 <- (<-answerResponse)

}