package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	log.Println("Connection on 127.0.0.1:2719 from: " + connection.LocalAddr().String())
	log.Println(input)

	var body string
	var status string
	input = strings.TrimSpace(input)
	switch input {
	case "/home":
		status = "200 OK"
		body = "Welcome to penguin net: " + connection.RemoteAddr().String() + "\n"
	default:
		status = "404 NOT FOUND"
		body = "unknown command: " + input + "\n"
	}
	response := makeResponse(status, body)
	connection.Write([]byte(response))
}

//standartizing response.. statuses need some structuring and the whole function requires some sanitizing before blindly appending things to the response i guess
func makeResponse(status string, body string) string {
      proto := "PENGO/0.1"
      contentLength := len(body)
      return proto + "\n" + status + "\n" +
              "Content-Length:" + strconv.Itoa(contentLength) + "\n\n" +
              body
}

func main() {
	listener, err := net.Listen("tcp", ":2719")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Listening on " + listener.Addr().String())
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}
