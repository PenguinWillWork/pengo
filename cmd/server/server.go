package main

import (
	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"pengo-proto/pengo"
)

func fileExists(path string, root *os.Root) bool {
	_, err := root.Stat("./" + path)
	return err == nil
}


func resolveExtension(input string, root *os.Root) (string, contentType string, err error) {
	if fileExists(input, root) {
		log.Println("file found as is returning: " + input)
		return input, filepath.Ext(input), nil
	}
	if fileExists(input + ".html", root) {
		log.Println("file comes with no ext -> returning html: " + input + ".html")
		return input + ".html", ".html", nil
	} 
	log.Println("no file:  " + input)
	return "", "text", errors.New("Input path's file doesn't exist\n");
}

func serverErrorResponse(connection net.Conn, err error) {
	response := pengo.MakeResponse(pengo.StatusInternalServerError, "text", []byte("Internal Server Error"))
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
		".html",
		notFoundContent,
	)
	log.Println(response)
	connection.Write([]byte(response))
	return;
}

func handleConnection(connection net.Conn, notFoundPath *string, root *os.Root) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	request, err := pengo.ParseRequest(reader) 
	if err != nil {
		return
	}

	var response []byte
	if request.RequestPath == "/" {
		request.RequestPath = "/index";
	}
	extensionResolvedInput, contentType, err := resolveExtension(request.RequestPath, root);
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
	response = pengo.MakeResponse(pengo.StatusOK, contentType, content)
	log.Println(contentType)
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
