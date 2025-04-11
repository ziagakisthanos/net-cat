package chat

import (
	"strings"
)

// HelpCommand checks if the provided command is a help command.
// If the command is "-h" or "-help", it sends the manual text to the client's output channel and returns true.
// Otherwise, it returns false so that other processing can continue.
func HelpCommand(command string, client *Client) bool {
	trimmed := strings.TrimSpace(command)
	if trimmed != "-h" && trimmed != "-help" {
		return false // not a help command, allow processing
	}

	// Send the manual directly to the client using the output channel.
	client.out <- GetManual()
	return true
}
