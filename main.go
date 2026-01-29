package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// database, global variable
var data = make(map[string]string)

// lock. 100 can read but only one can write
var mu sync.RWMutex

func main() {
	// Listen to TCP port (6379)
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379:", err)
		os.Exit(1)
	}

	fmt.Println("Lite-Redis listening on port 6379...")

	// Infinite loop to accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		// start a new goroutine for each connection to handle many requests concurrently
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// we read and cast to string
	buf := make([]byte, 1024) // 1KB buffer
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return // disconnected
		}
		// cast the data to string and print it (data is from a slice, so it's type is byte.)
		input := string(buf[:n])

		// "SET name archura \r\n" -> ["SET", "name", "archura"] we don't use string.Split cuz that might give us ""s
		parts := strings.Fields(input)
		// prevent PANIC if \r\n entered
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		switch command {

		case "PING":
			conn.Write([]byte("+PONG\r\n"))

		case "SET":
			if len(parts) != 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			key, value := parts[1], parts[2]
			// lock the mutex to ensure thread safety
			mu.Lock()
			data[key] = value
			mu.Unlock() // we good
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			// other people can read but can't write
			key := parts[1]
			mu.RLock()
			value, ok := data[key]
			mu.RUnlock()
			if !ok {
				conn.Write([]byte("$-1\r\n")) // Null Bulk String $-1, redis client doesn't recognize it as an error
			} else {
				resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
				conn.Write([]byte(resp))
			}
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
