package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"

	wsserver "github.com/AntonioMartinezFernandez/ws-performance-server/internal/ws_server"

	"golang.org/x/sys/unix"
)

var (
	addr           = flag.String("addr", ":8080", "server listen address")
	maxConns       = flag.Int64("maxconn", 2000000, "soft limit for connections (used as semaphore)")
	rlimit         = flag.Uint64("rlimit", 4000000, "set RLIMIT_NOFILE to this value")
	writeQueueSize = flag.Int("writequeuesize", 1024, "number of messages to queue per connection")
	pingInterval   = flag.Duration("pinginterval", 30*time.Second, "interval for websocket pings")

	activeConns *int64
)

func main() {
	flag.Parse()

	// Set GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU())

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Current active connections:", atomic.LoadInt64(activeConns))
		}
	}()

	// Increase RLIMIT_NOFILE (privileges needed if is higher than system limits)
	if err := setRlimit(*rlimit); err != nil {
		log.Printf("warning: no se pudo setear rlimit: %v", err)
	} else {
		log.Printf("rlimit nofile set to %d", *rlimit)
	}

	// Create listener with socket options
	listenerConfig := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return controlSocket(network, address, c)
		},
	}

	netListener, err := listenerConfig.Listen(
		context.Background(),
		"tcp",
		*addr,
	)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer netListener.Close()

	log.Printf("listening at %s (max connections %d)", *addr, *maxConns)

	server := &http.Server{
		Handler: wsserver.NewWebsocketHttpHandler(
			activeConns,
			*maxConns,
			*writeQueueSize,
			*pingInterval,
			wsserver.NewEchoAndPrinterMessageHandler(),
		),
	}

	// start serving on the created listener
	if err := server.Serve(netListener); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve: %v", err)
	}
}

func setRlimit(n uint64) error {
	rl := &unix.Rlimit{Cur: n, Max: n}
	return unix.Setrlimit(unix.RLIMIT_NOFILE, rl)
}

// control socket options (SO_REUSEADDR, SO_REUSEPORT)
func controlSocket(_, _ string, c syscall.RawConn) error {
	var ctrlErr error
	err := c.Control(func(fd uintptr) {
		// SO_REUSEADDR
		if e := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); e != nil {
			ctrlErr = e
			return
		}
		// SO_REUSEPORT (if available)
		_ = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	})
	if err != nil {
		return err
	}
	return ctrlErr
}
