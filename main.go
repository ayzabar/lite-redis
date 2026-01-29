package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Item struct {
	Value     string
	ExpiresAt int64
}

// database, global variable holds the item
var data = make(map[string]Item)

// lock. 100 can read but only one can write
var mu sync.RWMutex

func main() {
	// Listen to TCP port (6379)
	listener, err := net.Listen("tcp", ":6379")
	go startJanitor()
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
			if len(parts) < 3 {
				conn.Write([]byte("-ERR wrong number of arguments for 'set' command\r\n"))
				continue
			}
			key, value := parts[1], parts[2]

			var expiresAt int64 = 0

			// expiration control
			if len(parts) >= 5 && strings.ToUpper(parts[3]) == "EX" {
				seconds, err := strconv.Atoi(parts[4])
				if err == nil {
					// time + x secs
					expiresAt = time.Now().Add(time.Duration(seconds) * time.Second).UnixNano()
				}
			}

			// lock the mutex to ensure thread safety
			mu.Lock()
			data[key] = Item{
				Value:     value,
				ExpiresAt: expiresAt,
			}
			mu.Unlock() // we good
			conn.Write([]byte("+OK\r\n"))

		case "GET":
			if len(parts) != 2 {
				conn.Write([]byte("-ERR wrong number of arguments for 'get' command\r\n"))
				continue
			}
			key := parts[1]

			// FIX: used Lock instead of RLock
			// why? cuz we might delete during "Lazy Expiration"
			mu.Lock()
			item, ok := data[key]

			if !ok {
				mu.Unlock()
				conn.Write([]byte("$-1\r\n"))
				continue
			}

			// lazy expiration check
			if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
				delete(data, key)
				mu.Unlock()
				conn.Write([]byte("$-1\r\n"))
				continue
			}

			mu.Unlock() // we good
			resp := fmt.Sprintf("$%d\r\n%s\r\n", len(item.Value), item.Value)
			conn.Write([]byte(resp))

		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}
	}
}
func startJanitor() {
	for {
		time.Sleep(1 * time.Second)

		mu.Lock()

		now := time.Now().UnixNano()
		for key, item := range data {
			if item.ExpiresAt > 0 && now > item.ExpiresAt {
				delete(data, key)
			}
		}

		mu.Unlock()
	}
}
