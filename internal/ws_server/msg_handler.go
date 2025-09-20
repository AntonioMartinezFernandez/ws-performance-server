package wsserver

import (
	"fmt"
	"net"
)

type MessageHandler func(conn net.Conn, rcvMsg []byte)

func NewEchoAndPrinterMessageHandler() MessageHandler {
	return func(conn net.Conn, rcvMsg []byte) {
		// Echo back the received message
		if _, err := conn.Write(rcvMsg); err != nil {
			fmt.Println("echo write error:", err)
		}
		fmt.Printf("received message: %s", string(rcvMsg))
	}
}
