package chat

import (
	"net"
	"os"
	"sync"
)

// Client represents a connected chat user.
type Client struct {
	name string
	conn net.Conn
	out  chan string
}

// Server holds state for the chat server.
type Server struct {
	listenAddr string
	ln         net.Listener
	mu         sync.Mutex
	clients    map[string]*Client
	history    []string
	logFile    *os.File
}
