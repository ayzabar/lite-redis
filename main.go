package main

import (
	"fmt"
	"net"
	"os"
)

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

	// we read and print to the screen
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			return // disconnected
		}
		// cast the data to string and print it (data is from a slice, so it's type is byte.)
		fmt.Println(string("Received: " + string(buf[:n])))

		// answer (function accepts bytes, so we need to convert string to bytes. and ofc CRLF for the protocol)
		conn.Write([]byte("+PONG\r\n"))
	}
}
