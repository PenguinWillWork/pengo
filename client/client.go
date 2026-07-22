package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)


func makeRequest(ip string, request string) string {
	log.Println("requesting: " + ip + " request: " + request + "\n")
	connection, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatalln("Couldn't connect to the pengo server on " + ip + "\n");
	}
	connection.Write([]byte(request));
	response, err := io.ReadAll(connection);
	return string(response)
}

func parseInput(input string) (ip string, request string, err error){
	normalizedInput := strings.ToLower(input);
	proto := strings.Split(normalizedInput, "://")[0];
	inputWithoutProto := strings.Replace(normalizedInput, proto + "://", "", -1);
	if (proto != "pengo") {
		err := "Requested URI is not a Pengo URI";
		return "", "", errors.New(err); 
	}

	uriBody := strings.Split(inputWithoutProto, "/");
	return uriBody[0], "/" + uriBody[1], nil;
}

func main() {
	reader := bufio.NewReader(os.Stdin);
	fmt.Println("Enter pengo URI:");
	input, err := reader.ReadString('\n');
	if err != nil {
		log.Println(err);
	}

	ip, request, err := parseInput(input);
	if err != nil {
		log.Fatalln(err)
	}
	response := makeRequest(ip, request)
	fmt.Print(response)
}
