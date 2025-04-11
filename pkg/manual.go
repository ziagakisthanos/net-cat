package chat

import "strings"

// getManual constructs and returns the complete manual text.
// It uses ANSI escape sequences to clear the screen so that the help appears similar to an editor's manual mode.
func GetManual() string {
	var sb strings.Builder
	// Clear screen ANSI escape sequence (if supported by the client terminal)
	sb.WriteString("\033[H\033[2J")
	sb.WriteString("Netcat (TCP-Chat) Manual\n")
	sb.WriteString("--------------------------\n")
	sb.WriteString("Commands:\n")
	sb.WriteString("  -h, -help            Show this help manual.\n")
	sb.WriteString("  -name <new-name>     Change your nickname.\n")
	sb.WriteString("  <message>            Send a chat message to everyone.\n")
	sb.WriteString("\n")
	// sb.WriteString("Usage:\n")
	// sb.WriteString("	$ -name maria\n")
	// sb.WriteString("\n")
	// sb.WriteString("  This way, if your name was thanos, you can change it into maria.\n")
	// sb.WriteString("\n")
	// sb.WriteString("	$ -h or -help\n")
	sb.WriteString("\n")
	sb.WriteString("  To leave the server, press: \033[1mCtrl + C\033[0m\n")
	sb.WriteString("\n")
	sb.WriteString("Enjoy chatting!\n")
	return sb.String()
}
