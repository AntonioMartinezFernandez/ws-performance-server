package wsserver

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

func handleWebsocketConn(
	conn net.Conn,
	rcvMsgChan chan []byte,
	sendMsgChan chan []byte,
	doneChan chan struct{},
	activeConns *int64,
) {
	// set unique connection ID
	connectionID := uuid.New().String()
	log.Printf("new connection established: %s", connectionID)

	defer func() {
		conn.Close()                     // close connection on exit
		atomic.AddInt64(activeConns, -1) // decrement active connections counter
		close(doneChan)                  // stop writer goroutine
	}()

	// apply TCP_NODELAY
	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.SetNoDelay(true)
	}

	// writer goroutine
	go func() {
	writerFor:
		for {
			select {
			case msg := <-sendMsgChan:
				if err := wsutil.WriteServerMessage(conn, ws.OpText, msg); err != nil {
					fmt.Println("write error:", err)
					continue
				}
			case <-doneChan:
				break writerFor
			}
		}
	}()

	// reader loop
readerFor:
	for {
		// read message (non-blocking by wsutil design)
		rcvMessage, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {
			break readerFor
		}

		switch opCode {
		case ws.OpText, ws.OpBinary:
			// echo emulation with backpressure
			select {
			case sendMsgChan <- rcvMessage:
				// received message queued

			default:
				// full queue - close connection to prevent wear out
				log.Printf("write queue full - closing connection")
				break readerFor
			}

		case ws.OpPing:
			// respond to ping
			if err := wsutil.WriteServerMessage(conn, ws.OpPong, rcvMessage); err != nil {
				fmt.Println("pong error:", err)
				continue
			}
			fmt.Println("received ping from client")

		case ws.OpPong:
			fmt.Println("received pong from client")

		case ws.OpClose:
			fmt.Println("received close frame")
			break readerFor
		}
	}
}
