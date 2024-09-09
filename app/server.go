package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	ArrayPrefix = '*'
	BulkPrefix  = '$'
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
	defer conn.Close()

	for {
		// Create a buffer to read
		buf := make([]byte, 128)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// Parse the command from the buffer
		command, _ := parseCommand(buf[:n])

		switch strings.ToLower(command[0]) {
			case "ping":
				conn.Write([]byte("+PONG\r\n"))
			case "echo":
				if len(command) > 1 {
					args := strings.Join(command[1:], " ")
					// RESP bulk string format: $<length>\r\n<arg>\r\n
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


// Parse RESP protocol format
func parseCommand(buf []byte) ([]string, error) {
	// Go through each byte in the buffer
	i := 0
	if i >= len(buf) {
		return nil, fmt.Errorf("invalid command")
	}
	if buf[i] != ArrayPrefix {
		return nil, fmt.Errorf("invalid command format")
	}
	i++
	var length int
	for i < len(buf) && buf[i] >= '0' && buf[i] <= '9' {
		// Convert the character to integer
		length = length*10 + int(buf[i]-'0')
		i++
	}
	i, err := expect(buf, i, "\r\n")
	if err != nil {
		return nil, err
	}

	var args []string
	var arg string
	for j := 0; j < length; j++ {
		arg, i, err = parseBulkString(buf, i)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}

// $4\r\nECHO\r\n
func parseBulkString(buf []byte, i int) (string, int, error) {
	// Check if the buffer contains '$' in the beginning
	if i >= len(buf) || buf[i] != BulkPrefix {
		return "", i, errors.New("Expecting $")
	}
	i++
	var length int
	for i < len(buf) && buf[i] >= '0' && buf[i] <= '9' {
		// Convert the character to integer
		length = length*10 + int(buf[i]-'0')
		i++
	}
	i, err := expect(buf, i, "\r\n")
	if err != nil {
		return "", i, err
	}
	// Extract the string
	if i+length > len(buf) {
        return "", i, fmt.Errorf("buffer too short for the expected length")
    }
    bulkString := string(buf[i : i+length])
    i += length

    i, err = expect(buf, i, "\r\n")
    if err != nil {
        return "", i, err
    }
	return bulkString, i, nil
}

// Check if the buffer contains '\r\n' at the specified index
func expect(buf []byte, i int, exp string) (int, error) {
	if i+len(exp) <= len(buf) && string(buf[i:i+len(exp)]) == exp {
		return i + len(exp), nil
	}
	return i, errors.New("Expecting " + exp)
}
