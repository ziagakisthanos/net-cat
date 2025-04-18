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

// NameChange processes a nickname change command.
// If the received message starts with "-namechange ", it will attempt to change the client's nickname.
// The parameter currentName is a pointer to the variable holding the client's current name.
// Returns true if the message was processed as a nickname change (even if the change failed), false otherwise.
func NameCommand(command string, s *Server, client *Client, currentName *string) bool {
	// Check if the message begins with the command "/nick "
	if (strings.HasPrefix(command, "-n")) || (strings.HasPrefix(command, "-name")) {

		tokens := strings.Fields(command)
		if len(tokens) < 2 {
			if tokens[0] == "-n" {
				client.out <- NamePrompt()
				return true
			}
			return false
		}
		// Extract the new name token from the command.
		newName := tokens[1]
		//prompt the user on blank new name
		if newName == "" {
			client.out <- NamePrompt()
			return true
		}

		//if new name is empty, prompt the user
		if !isMessageValid([]byte(newName)) {
			client.out <- "[SERVER]: Invalid new name. Please try again."
			return true
		}

		// Lock the server state while checking and updating the client name.
		s.mu.Lock()
		// Check if the new name is already taken.
		if _, exists := s.clients[newName]; exists {
			s.mu.Unlock()
			client.out <- "[SERVER]: Name already taken. Please choose a different name."
			return true
		}

		// Process the name change.
		oldName := *currentName
		// Remove the client from the mapping keyed by the old name.
		delete(s.clients, oldName)
		// Update the client's stored name.
		client.name = newName
		// Add the client back to the mapping with the new name.
		s.clients[newName] = client
		// Update the currentName variable to reflect the new nickname.
		*currentName = newName
		s.mu.Unlock()

		// Create a notification message.
		nameChangeMsg := fmt.Sprintf("[%s][SERVER]: %s changed their name to %s",
			time.Now().Format("2006-01-02 15:04:05"),
			oldName,
			newName)

		// Add the nickname change message to history and broadcast it to all clients.
		s.addHistory(nameChangeMsg)
		s.broadcast(nameChangeMsg, "")
		return true
	}
	return false // not a name command; let the caller process it as a normal message.
}
