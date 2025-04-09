package chat

import "fmt"

// getHistory returns a copy of the current in-memory chat history.
func (s *Server) getHistory() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	histCopy := make([]string, len(s.history))
	copy(histCopy, s.history)
	return histCopy
}

// addHistory appends a message to the in-memory chat history.
func (s *Server) addHistory(message string) {
	s.mu.Lock()
	s.history = append(s.history, message)
	s.mu.Unlock()
	// outputs message on server side
	fmt.Println(message)
}
