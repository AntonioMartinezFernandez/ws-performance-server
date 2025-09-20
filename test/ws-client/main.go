package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var (
	serverAddr  = "localhost:8080"
	message     = `{"foo":"bar %s"}`
	msgInterval = 1500 * time.Millisecond
	clients     = 2
)

func main() {
	// Start 10 clients
	for i := range clients {
		go connectClient(i + 1)
	}

	// Run indefinitely
	select {}
}

func connectClient(clientNumber int) {
	// Connect to the server
	c, _, err := websocket.DefaultDialer.Dial("ws://"+serverAddr, nil)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer c.Close()

	// Receive and print messages
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				return
			}
			log.Printf("recv: %s", msg)
		}
	}()

	// Send messages periodically
	ticker := time.NewTicker(msgInterval)
	defer ticker.Stop()

	for range ticker.C {
		msg := fmt.Sprintf(message, strconv.Itoa(clientNumber))
		err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("write error: %v", err)
			return
		}
		log.Printf("sent: %s", msg)
	}
}
