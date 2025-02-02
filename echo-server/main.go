package main

import (
	"fmt"
	"net"
	"os"
)

// StartEchoServer it takes a port and listen to it, takes the input and buffer it
// then echo it back to the server
func StartEchoServer(port int) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Println("Echo server listening on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

// handleConnection takes a TCP conn object and handle reading and writing
func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Connection closed or error reading:", err)
			return
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("Connection closed or writing reading:", err)
			return
		}
	}
}

func main() {

}
