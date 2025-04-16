package chat

import "strings"

// isMessageValid checks that the provided data (name or message) contains at least one non-whitespace character.
func isMessageValid(data []byte) bool {
	trimmed := strings.TrimSpace(string(data))
	return len(trimmed) != 0
}

// TODO:
// 		FIX	client message spawn twice
//
