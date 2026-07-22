package main

import (
	"io"
	"log"
	"net"
)

func handleConnection(connection net.Conn) {
	log.Println(connection.RemoteAddr(), ": Says hello")
	io.Copy(connection, connection)
	connection.Close()
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1927")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Listening on 127.0.0.1:1927");
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
