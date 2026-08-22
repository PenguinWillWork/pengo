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

//resolve file extension for request path, if no extension is found, try to resolve with .html extension
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

func requestMalformedResponse(connection net.Conn, err error) {
	response := pengo.MakeResponse(pengo.StatusBadRequest, "text", []byte("Bad Request"))
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
	connection.Write(response)
	return;
}

func handleFetchRequest(connection net.Conn, notFoundPath *string, root *os.Root, request *pengo.Request) ([]byte, error) {
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
		return []byte{}, errors.New("File not found")
	}
    if err != nil {
		serverErrorResponse(connection, err);
		return []byte{}, err;
	}
	response := pengo.MakeResponse(pengo.StatusOK, contentType, content)
	return response, nil;
}


func handleConnection(connection net.Conn, notFoundPath *string, root *os.Root) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	request, err := pengo.ParseRequest(reader)

	if err != nil || request.Method != "fetch" && request.Method != "submit" {
		log.Println("Error parsing request: " + err.Error())
		requestMalformedResponse(connection, err)
		return;
	}

	var response []byte
	switch request.Method {
		case "fetch":
			response, err = handleFetchRequest(connection, notFoundPath, root, &request)
			if err != nil {
				log.Println("Error handling fetch request: " + err.Error())
				return;
			}
		case "submit":
			log.Println("Submit method not implemented yet")		
	}
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
