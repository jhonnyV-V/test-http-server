package main

import (
	"fmt"
	"net"
	"os"
)

const (
	OK = "HTTP/1.1 200 OK\r\n\r\n"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	connection.Write([]byte(OK))
	connection.Close()
}
