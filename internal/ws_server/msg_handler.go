package wsserver

import (
	"fmt"
)

type MessageHandler func(rcvMsg []byte, sendMsgChan chan<- []byte)

func NewEchoAndPrinterMessageHandler() MessageHandler {
	return func(rcvMsg []byte, sendMsgChan chan<- []byte) {
		// echo the received message
		sendMsgChan <- rcvMsg

		// print the received message to the console
		fmt.Printf("received message: %s\n", string(rcvMsg))
	}
}
