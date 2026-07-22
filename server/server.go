package main

import (
	"bufio"
	"log"
	"net"
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
	response := "PENGO/0.1 " + status + "\n" + body
	connection.Write([]byte(response))
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
