package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Server is listening on port 6379...")

	// Keep accepting incoming connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			os.Exit(1)
		}

		// Using goroutine to handle each client connection separately
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	for {
		// Create a buffer to read the client's command
		buf := make([]byte, 128)
		conn.Read(buf)
		conn.Write([]byte("+PONG\r\n"))
	}
}
