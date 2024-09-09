package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

var storage = make(map[string]ExpiringValue)
type ExpiringValue struct {
	Value      string
	ExpireAt   time.Time
}

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
			if len(command) >= 3 {
				key := command[1]
				value := command[2]
				expireTime := time.Time{}
				if len(command) > 3 {
					if strings.ToLower(command[3]) == "px" && len(command) > 4 {
						duration, err := time.ParseDuration(command[4] + "ms")
						if err != nil {
							conn.Write([]byte("-ERR invalid PX duration\r\n"))
							return
						}
						expireTime = time.Now().Add(duration)
					}
				}
				storage[key] = ExpiringValue{Value: value, ExpireAt: expireTime}
				conn.Write([]byte("+OK\r\n"))
			} else {
				conn.Write([]byte("-ERR SET command requires a key and a value\r\n"))
			}
		// get a value by key
		case "get":
			if len(command) == 2 {
				key := command[1]
				value, ok := storage[key]
				if ok && (value.ExpireAt.IsZero() || value.ExpireAt.After(time.Now())) {
					response := fmt.Sprintf("$%d\r\n%s\r\n", len(value.Value), value.Value)
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
