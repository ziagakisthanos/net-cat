package main

import (
	"fmt"
	"log"
	chat "netcat/pkg"
	"os"
)

func main() {
	// Allowed usage: no command-line argument (default port) or one argument specifying the port.
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}

	port := ""
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		port = "8989" // Default port
	}
	addr := ":" + port

	server := chat.NewServer(addr)
	fmt.Println("Server online...")
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
