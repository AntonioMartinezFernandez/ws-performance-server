# Websockets Performance Server

A high performance Websockets server written in Go. It is designed to handle a large number of concurrent WebSocket connections efficiently, making it suitable for applications that require real-time communication, such as chat applications, live updates, and multiplayer games.

## Features

- Efficient handling of WebSocket connections using the `gobwas/ws` library.
- Configurable maximum number of concurrent connections.
- Configurable write queue size per connection.
- Configurable ping interval to keep connections alive.
- Graceful shutdown on receiving termination signals.
