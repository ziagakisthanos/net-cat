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

	// Write the log into the file if it exists.
    	if s.logFile != nil {
    		_, err := fmt.Fprintln(s.logFile, message)
    		if err != nil {
    			// If writing to the file fails, log the error to the console.
    			fmt.Println("Error writing to log file:", err)
    		}
    	}
}
