package main

import (
	"fmt"
	"net"
	"strings"
)

var storage = make(map[string]string)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 128)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		command, _ := parseCommand(buf[:n])

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
		// set a key to a value
		case "set":
			if len(command) == 3 {
				key := command[1]
				value := command[2]
				storage[key] = value
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR SET command requires a key and a value\r\n"))
			}
		// get a value by key
		case "get":
			if len(command) == 2 {
				key := command[1]
				value, ok := storage[key]
				if ok {
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
					conn.Write([]byte(response))
				} else {
					conn.Write([]byte("$-1\r\n"))
				}
			} else {
				conn.Write([]byte("-ERR GET command requires a key\r\n"))
			}
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
