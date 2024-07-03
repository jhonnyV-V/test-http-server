package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	OK        = "HTTP/1.1 200 OK\r\n\r\n"
	NOT_FOUND = "HTTP/1.1 404 Not Found\r\n\r\n"
)

var FileDir string

func isCRLF(x string) bool {
	crlf := []byte{13, 10}
	inBytes := []byte(x)

	if len(inBytes) != 2 {
		return false
	}
	return inBytes[0] == crlf[0] && inBytes[1] == crlf[1]
}

func okWithBody(body, ctype string) []byte {
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s",
		ctype,
		len(body),
		body,
	)
	fmt.Println("response", response)
	return []byte(response)
}

func parseHeaders(buff *bufio.Reader) []string {
	var headers []string
	header := "x"
	for !isCRLF(header) {
		header, _ = buff.ReadString('\n')
		headers = append(headers, strings.TrimSpace(header))
	}

	return headers
}

func echo(fullpath string) []byte {
	path := strings.Split(fullpath, "/")
	message := path[2]

	return okWithBody(message, "text/plain")
}

func userAgent(headers []string) []byte {
	agent := ""
	for _, v := range headers {
		if strings.HasPrefix(v, "User-Agent") {
			agent = strings.TrimSpace(strings.Split(v, ":")[1])
		}
	}

	return okWithBody(agent, "text/plain")
}

func getFiles(reqPath string) []byte {
	path := strings.Split(reqPath, "/")
	fileName := path[2]

	raw, err := os.ReadFile(FileDir + fileName)
	if err != nil {
		fmt.Println("error: ", err)
		fmt.Println("fileName: ", fileName)
		fmt.Println("fileDir: ", FileDir)
		return []byte(NOT_FOUND)
	}

	return okWithBody(string(raw), "application/octet-stream")
}

func HandleConnection(connection net.Conn) {
	defer connection.Close()
	buff := bufio.NewReader(connection)
	raw, err := buff.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(raw)

	args := strings.Split(raw, " ")
	if args[1] == "/" {
		_, err = connection.Write([]byte(OK))
	} else if strings.HasPrefix(args[1], "/echo/") {
		_, err = connection.Write(echo(args[1]))
	} else if strings.HasPrefix(args[1], "/user-agent") {
		headers := parseHeaders(buff)
		response := userAgent(headers)
		_, err = connection.Write(response)
	} else if strings.HasPrefix(args[1], "/files") {
		_, err = connection.Write(getFiles(args[1]))
	} else {
		_, err = connection.Write([]byte(NOT_FOUND))
	}

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.StringVar(&FileDir, "directory", "", "--directory /tpm/something")
	flag.Parse()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleConnection(connection)
	}
}
