package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	OK        = "HTTP/1.1 200 OK\r\n\r\n"
	NOT_FOUND = "HTTP/1.1 404 Not Found\r\n\r\n"
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

	buff := bufio.NewReader(connection)
	raw, err := buff.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(raw)

	args := strings.Split(raw, " ")
	if args[1] == "/" {
		_, err = connection.Write([]byte(OK))
	} else {
		_, err = connection.Write([]byte(NOT_FOUND))
	}

	if err != nil {
		fmt.Println(err)
	}
	connection.Close()
}
