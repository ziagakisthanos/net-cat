package chat

import "strings"

// isMessageValid checks that the provided data (name or message) contains at least one non-whitespace character.
func isMessageValid(data []byte) bool {
	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return false
	}
	if strings.HasPrefix(string(data), "-name ") {
		return false
	}
	return true
}

// TODO:
// 		FIX -name tag
//				client message spawn twice
//					validation refactor can help ?
//		ADD	private messages
//		DELETE exclude feature
//
