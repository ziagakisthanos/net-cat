package chat

import (
	"fmt"
	"strings"
	"time"
)

func NamePrompt() string {
	var sb strings.Builder
	// ANSI escape sequence to clear the screen.
	sb.WriteString("\n")
	sb.WriteString("Usage:\n")
	sb.WriteString("-name <new-name>\nor -n <new-name>")
	return sb.String()
}

// NameCommand processes a nickname-change command.
// Returns true if the line was a name-change request (even if it failed), false otherwise.
func NameCommand(command string, s *Server, client *Client, currentName *string) bool {
    // Split the input into whitespace-separated tokens.
    tokens := strings.Fields(command)
    if len(tokens) == 0 {
        return false
    }

    // Check if the first token is exactly "-n" or "-name"
    cmd := tokens[0]
    if cmd != "-n" && cmd != "-name" {
        // Not a name-change command; let the caller treat it as a normal message.
        return false
    }

    // Must have a second token for the new name
    if len(tokens) < 2 {
        client.out <- NamePrompt()
        return true
    }
    newName := tokens[1]

    // Validate the new name
    if !isMessageValid([]byte(newName)) {
        client.out <- "[SERVER]: Invalid new name. Please try again."
        return true
    }

    // Attempt to change the name under lock
    s.mu.Lock()
    if _, exists := s.clients[newName]; exists {
        s.mu.Unlock()
        client.out <- "[SERVER]: Name already taken. Please choose a different name."
        return true
    }
    oldName := *currentName
    delete(s.clients, oldName)
    client.name = newName
    s.clients[newName] = client
    *currentName = newName
    s.mu.Unlock()

    // Broadcast the name change
    msg := fmt.Sprintf("[%s][SERVER]: %s changed their name to %s",
        time.Now().Format("2006-01-02 15:04:05"),
        oldName, newName)
    s.addHistory(msg)
    s.broadcast(msg, "")

    return true
}