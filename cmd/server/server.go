package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"pengo-proto/pengo"
	"strings"
)


func handleConnection(connection net.Conn, notFoundPath *string, root *os.Root) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	var response string
	input = strings.TrimSpace(input)
    content, err := root.ReadFile("./" + input + ".html")

	if len(content) == 0 {
		default404 := "<!doctype html><html><body><h1>404 Not Found</h1></body></html>"

		notFoundContent := []byte(default404)

		if notFoundPath != nil {
			content, err := root.ReadFile("./" + *notFoundPath + ".html")


			if err == nil {
				notFoundContent = content
			}
		}

		response = pengo.MakeResponse(
			pengo.StatusNotFound,
			string(notFoundContent),
		)
		connection.Write([]byte(response))
		return;
	}
    if err != nil {
        response = pengo.MakeResponse(pengo.StatusInternalServerError, err.Error())
		connection.Write([]byte(response))
		return;
	}

	response = pengo.MakeResponse(pengo.StatusOK, string(content))
	connection.Write([]byte(response))
}

func main() {
	notFoundPath := "/404";
	root, err := os.OpenRoot("./")
	
	if err != nil {
		log.Fatal(err)
	}
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

		go handleConnection(conn, &notFoundPath, root)
	}

}
