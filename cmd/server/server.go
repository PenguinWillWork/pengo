package main

import (
	"bufio"
	"errors"
	"flag"
	"log"
	"mime"
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
		return input, mime.TypeByExtension(filepath.Ext(input)), nil
	}
	if fileExists(input + ".html", root) {
		return input + ".html", mime.TypeByExtension(".html"), nil
	} 
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
		log.Println("Error parsing request: " + err.Error())
		serverErrorResponse(connection, err)
		return;
	}

	log.Println("Received request: " + request.Method + " " + request.RequestPath + " " + request.Host)
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
	// log.Println(contentType)
	connection.Write([]byte(response))
}

func main() {
	port := flag.String("port", "2719", "port to listen on")
	rootDir := flag.String("root", ".", "directory to serve")
	flag.Parse()

	notFoundPagePath := "/404.html";
	
	root, err := os.OpenRoot(*rootDir)          
	
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", ":"+*port)   
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
