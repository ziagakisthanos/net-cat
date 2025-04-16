package chat

import (
	"fmt"
	"strings"
	"time"
)

func WhisperPrompt() string {
	var sb strings.Builder
	// ANSI escape sequence to clear the screen.
	sb.WriteString("\n")
	sb.WriteString("Usage:\n")
	sb.WriteString("-w [recipient] [message]\n")
	return sb.String()
}

// WhisperCommand processes a private message (whisper) command.
// The command parameter should contain the full text (e.g. "-w abby hello there").
// The function sends the message privately to the specified recipient using the server instance s.
// The sender's name is provided via senderName.
// Returns true if the command was recognized and processed, false otherwise.
func WhisperCommand(command string, s *Server, client *Client, senderName string) bool {
	if strings.HasPrefix(command, "-w") || strings.HasPrefix(command, "-whisper") {
		// Tokenize the command into its constituent parts.
		tokens := strings.Fields(command)
		if len(tokens) < 3 {
			if tokens[0] == "-w" {
				client.out <- WhisperPrompt()
				return true
			}
			return false
		}

		// The second token is the recipient's name.
		target := tokens[1]

		//check if the client exists
		s.mu.Lock()
		//map lookup
		//value, ok := map[key]
		_, targetExists := s.clients[target]
		s.mu.Unlock()
		if !targetExists {
			client.out <- fmt.Sprintf("%s is not online or does not exist.", target)
			return true
		}
		// The remainder of the tokens form the message body.
		msg := strings.Join(tokens[2:], " ")

		// Format the whisper message with a timestamp and indicate its a private message.
		formattedMessage := fmt.Sprintf("[%s][%s] whispers: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			senderName,
			msg)

		// Use the server's broadcast function to send the message only to the specified target.
		s.broadcast(formattedMessage, target)
		return true
	}
	return false
}
