package main

import (
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 128)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		command, _ := parseCommand(buf[:n])  // Assuming this function is defined in parser.go

		switch strings.ToLower(command[0]) {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			if len(command) > 1 {
				args := strings.Join(command[1:], " ")
				response := fmt.Sprintf("$%d\r\n%s\r\n", len(args), args)
				_, err = conn.Write([]byte(response))
				if err != nil {
					fmt.Println("Error writing response:", err)
					return
				}
			} else {
				conn.Write([]byte("-ERR ECHO command requires an argument\r\n"))
			}
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
