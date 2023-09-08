package net

import (
	"testing"
)

func TestInitOffer(t *testing.T) {
	
	offerChan := make(chan string)
	answerResponseEncoded := make(chan string)

	go InitOffer(offerChan, answerResponseEncoded)

	t.Log(<-offerChan)

}