package chat

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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
	mu         sync.Mutex
	clients    map[string]*Client
	history    []string
	logFile    *os.File
}

// NewServer creates a new Server instance using the provided address.
func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		clients:    make(map[string]*Client),
		history:    make([]string, 0),
		logFile:    SetupLogFile(listenAddr),
	}
}

// Start opens the listening socket and continuously accepts incoming connections.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	// Inform server console about the port.
	fmt.Printf("Listening on the port :%s\n", s.listenAddr[1:]) // remove leading ":"

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		fmt.Println("New connection: ", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

// server.go (only the high-level flow; helper methods below)
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// 1️⃣ Register (handshake + name)
	client, name, err := s.registerClient(conn, reader)
	if err != nil {
		return
	}
	defer s.deregisterClient(name, client)

	// 2️⃣ Send existing history
	s.sendHistory(client)

	// 3️⃣ Announce join
	s.announceJoin(name)

	// 4️⃣ Start writer
	go s.clientWriter(client)

	// 5️⃣ Enter main read+broadcast loop
	s.handleMessages(reader, client, &name)

	// 6️⃣ Announce leave
	s.announceLeave(name)
}

// handleMessages is the core per-client read loop.
func (s *Server) handleMessages(reader *bufio.Reader, client *Client, name *string) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Printf("Read error from %s: %v\n", *name, err)
			}
			return
		}

		// erase raw echo (ANSI only—make this optional/configurable)
		client.conn.Write([]byte("\x1b[1A\r\x1b[K"))

		text := strings.TrimSpace(line)
		if !isMessageValid([]byte(text)) {
			continue
		}

		// commands: whisper, help, name-change
		if strings.HasPrefix(text, "-") {
			switch {
			case WhisperCommand(text, s, client, *name):
			case HelpCommand(text, client):
			case NameCommand(text, s, client, name):
			default:
				client.out <- text + " is not a command. try -h"
			}
			continue
		}

		// normal broadcast
		msg := fmt.Sprintf("[%s][%s]: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			*name,
			text)
		s.addHistory(msg)
		s.broadcast(msg, "")
	}
}

// registerClient does the name handshake and client-map insertion.
func (s *Server) registerClient(conn net.Conn, reader *bufio.Reader) (*Client, string, error) {
	conn.Write([]byte(welcomeBanner))

	for {
		nameLine, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading name:", err)
			return nil, "", err
		}
		name := strings.TrimSpace(nameLine)
		if !isMessageValid([]byte(name)) {
			conn.Write([]byte("Invalid name. Please try again.\n"))
			continue
		}

		s.mu.Lock()
		if len(s.clients) >= maxClients {
			s.mu.Unlock()
			conn.Write([]byte("Server full. Try later.\n"))
			return nil, "", fmt.Errorf("server full")
		}
		if _, exists := s.clients[name]; exists {
			s.mu.Unlock()
			conn.Write([]byte("Name taken. Choose another.\n"))
			continue
		}

		client := &Client{name: name, conn: conn, out: make(chan string, 10)}
		s.clients[name] = client
		s.mu.Unlock()
		return client, name, nil
	}
}

// broadcast sends a message to clients.
// If target is empty, it broadcasts to all connected clients.
// If target is non-empty, it sends the message only to the client with that name (private message).
func (s *Server) broadcast(message, target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	//private message to target if specified
	if target != "" {
		if client, ok := s.clients[target]; ok {
			client.out <- message
		}
		return
	}
	//broadcast to all clients.
	for _, client := range s.clients {
		client.out <- message
	}
}

// deregisterClient removes the client from the map and closes its out-channel.
func (s *Server) deregisterClient(name string, client *Client) {
	s.mu.Lock()
	delete(s.clients, name)
	s.mu.Unlock()
	close(client.out)
}

func (s *Server) announceJoin(name string) {
	joinMsg := fmt.Sprintf("[%s][SERVER]: %s joined our chat",
		time.Now().Format("2006-01-02 15:04:05"), name)
	s.addHistory(joinMsg)
	s.broadcast(joinMsg, "")
}

func (s *Server) announceLeave(name string) {
	leaveMsg := fmt.Sprintf("[%s][SERVER]: %s left our chat",
		time.Now().Format("2006-01-02 15:04:05"), name)
	s.addHistory(leaveMsg)
	s.broadcast(leaveMsg, "")
}

// sendHistory streams the in-memory history down to this client.
func (s *Server) sendHistory(client *Client) {
	for _, msg := range s.getHistory() {
		client.out <- msg
	}
}
