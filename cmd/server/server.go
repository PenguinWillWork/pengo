package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"pengo-proto/pengo"
	"strings"
)

func fileExists(path string, root *os.Root) bool {
	_, err := root.Stat("./" + path)
	log.Println(err)
	return err == nil
}


func resolveExtension(input string, root *os.Root) (string, error) {
	if fileExists(input, root) {
		log.Println("file found as is returning: " + input)
		return input, nil
	}
	if fileExists(input + ".html", root) {
		log.Println("file comes with no ext -> returning html: " + input + ".html")
		return input + ".html", nil
	} 
	log.Println("no file:  " + input)
	return "", errors.New("Input path's file doesn't exist\n");
}

func serverErrorResponse(connection net.Conn, err error) {
	response := pengo.MakeResponse(pengo.StatusInternalServerError, "Internal Server Error")
	connection.Write([]byte(response))
	return;
}

func notFoundResponse(connection net.Conn, notFoundPath *string, root *os.Root) {
	default404 := "<!doctype html><html><body><h1>404 Not Found</h1></body></html>"

	notFoundContent := []byte(default404)

	if notFoundPath != nil {
		content, err := root.ReadFile("./" + *notFoundPath)
		if err == nil {
			notFoundContent = content
		}
	}

	response := pengo.MakeResponse(
		pengo.StatusNotFound,
		string(notFoundContent),
	)
	log.Println(response)
	connection.Write([]byte(response))
	return;
}

func handleConnection(connection net.Conn, notFoundPath *string, root *os.Root) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	var response string
	input = strings.TrimSpace(input)

	if input == "/" {
		input = "/index";
	}
	extensionResolvedInput, err := resolveExtension(input, root);
	if err != nil {
		serverErrorResponse(connection, err);
	}
    content, err := root.ReadFile("./" + extensionResolvedInput)

	if len(content) == 0 {
		notFoundResponse(connection, notFoundPath, root)
		return;
	}
    if err != nil {
		serverErrorResponse(connection, err);
		return;
	}

	response = pengo.MakeResponse(pengo.StatusOK, string(content))
	connection.Write([]byte(response))
}

func main() {
	notFoundPagePath := "/404.html";
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

		go handleConnection(conn, &notFoundPagePath, root)
	}

}
