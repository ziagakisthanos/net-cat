package chat

import (
	"fmt"
	"os"
)

// SetupLogFile opens (or creates) the log file in append mode based on the server's listen address.
// It differentiates the log file based on the port. For example, if listenAddr is ":8989", the file will be "log8989".
func SetupLogFile(listenAddr string) *os.File {
	// Extract the port from listenAddr (assuming format ":port").
	port := listenAddr
	if len(port) > 0 && port[0] == ':' {
		port = port[1:]
	}
	fileName := "logs/log" + port // e.g., log8989 for port 8989

	// Open the file in append mode (O_APPEND), create it if it does not exist (O_CREATE),
	// and open it for writing (O_WRONLY). Mode is set to 0644.
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file", fileName, ":", err)
		return nil
	}
	return logFile
}
