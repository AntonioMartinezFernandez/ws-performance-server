package wsserver

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
)

// http handler for creating websockets connections using gobwas/ws
func NewWebsocketConnectionHttpHandler(
	activeConns *int64,
	maxConns int64,
	writeQueueSize int,
	pingInterval time.Duration,
	msgHandler MessageHandler,
) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// reject connections if over soft limit
		if atomic.LoadInt64(activeConns) >= maxConns {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("server overloaded"))
			return
		}

		// upgrade connection
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("upgrade failed: %v", err)
			return
		}

		// increment active connections counter
		atomic.AddInt64(activeConns, 1)

		// buffered channel as write queue for each connection
		sendMsgChan := make(chan []byte, writeQueueSize)

		// doneChan channel to signal writer goroutine to exit
		doneChan := make(chan struct{})

		// start ping and connection handlers
		go handleWebsocketPing(conn, doneChan, pingInterval)
		go handleWebsocketConn(conn, sendMsgChan, doneChan, activeConns, msgHandler)
	}

	return http.HandlerFunc(handler)
}
