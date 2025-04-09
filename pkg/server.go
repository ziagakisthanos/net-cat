package chat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

// NewServer creates a new Server instance using the provided address.
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		clients:    make(map[string]*Client),
		history:    make([]string, 0),
	}
}

// Start opens the listening socket and continuously accepts incoming connections.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	// Inform server console about the port.
	fmt.Printf("Listening on the port :%s\n", s.listenAddr[1:]) // remove leading ":"

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		fmt.Println("New connection: ", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

// handleConnection manages the handshake and communication with a client.
func (s *Server) handleConnection(conn net.Conn) {
	// Use bufio to read from the connection.
	reader := bufio.NewReader(conn)
	// Send welcome banner (Linux logo + prompt for name).
	conn.Write([]byte(welcomeBanner))

	// Read the client's name.
	nameLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading name:", err)
		conn.Close()
		return
	}
	name := strings.TrimSpace(nameLine)
	if !isMessageValid([]byte(name)) {
		conn.Write([]byte("Invalid name. Disconnecting.\n"))
		conn.Close()
		return
	}

	// Lock the server state to check connection limits and uniqueness.
	s.mu.Lock()
	if len(s.clients) >= maxClients {
		s.mu.Unlock()
		conn.Write([]byte("Server full. Please try again later.\n"))
		conn.Close()
		return
	}
	if _, exists := s.clients[name]; exists {
		s.mu.Unlock()
		conn.Write([]byte("Name already taken. Disconnecting.\n"))
		conn.Close()
		return
	}
	// New client
	client := &Client{
		name: name,
		conn: conn,
		out:  make(chan string, 10),
	}
	s.clients[name] = client
	s.mu.Unlock()

	// Send the full in-memory chat history to the new client.
	for _, msg := range s.getHistory() {
		client.out <- msg
	}

	// Broadcast that a new client has joined.
	joinMsg := fmt.Sprintf("[%s][SERVER]: %s joined our chat...", time.Now().Format("2006-01-02 15:04:05"), name)
	s.addHistory(joinMsg)
	s.broadcast(joinMsg, "")

	// Start the writer goroutine for this client.
	go s.clientWriter(client)

	// Read loop
	for {
		msgLine, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("Read error from %s: %v\n", name, err)
			}
			break
		}
		if !isMessageValid([]byte(msgLine)) {
			// Do not process empty messages.
			continue
		}
		trimmedMsg := strings.TrimSpace(msgLine)
		formattedMsg := fmt.Sprintf("[%s][%s]: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			name,
			trimmedMsg)
		s.addHistory(formattedMsg)
		s.broadcast(formattedMsg, "")
	}

	// Client disconnects.
	s.mu.Lock()
	delete(s.clients, name)
	s.mu.Unlock()
	leaveMsg := fmt.Sprintf("[%s][SERVER]: %s left our chat...", time.Now().Format("2006-01-02 15:04:05"), name)
	s.addHistory(leaveMsg)
	s.broadcast(leaveMsg, "")
	conn.Close()
}

// broadcast sends message to all connected clients except (optionally) the one specified by 'exclude'.
func (s *Server) broadcast(message, exclude string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for clientName, client := range s.clients {
		if clientName == exclude {
			continue
		}
		select {
		case client.out <- message:
		default:
			// Skip if the client's outbound channel is full.
		}
	}
}
