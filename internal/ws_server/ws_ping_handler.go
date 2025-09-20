package wsserver

import (
	"fmt"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func handleWebsocketPing(conn net.Conn, doneChan chan struct{}, pingInterval time.Duration) {
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

pingFor:
	for {
		select {
		case <-ticker.C:
			if err := wsutil.WriteServerMessage(conn, ws.OpPing, []byte("_ping_")); err != nil {
				fmt.Println("ping error:", err)
				continue
			}

		case <-doneChan:
			break pingFor
		}
	}
}
