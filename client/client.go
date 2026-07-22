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
	connection.Write([]byte(request + "\n"));
	response, err := io.ReadAll(connection);
	return string(response)
}

func parseInput(input string) (ip string, request string, err error){
	normalizedInput := strings.TrimSpace(strings.ToLower(input));
	proto := strings.Split(normalizedInput, "://")[0];
	inputWithoutProto := strings.Replace(normalizedInput, proto + "://", "", -1);
	if (proto != "pengo") {
		err := "Requested URI is not a Pengo URI";
		return "", "", errors.New(err);
	}

	uriBody := strings.SplitN(inputWithoutProto, "/", 2);
	host := uriBody[0];
	path := "/home";
	if len(uriBody) > 1 && uriBody[1] != "" {
		path = "/" + uriBody[1];
	}
	return host, path, nil;
}

func resolveDns(host string) (resolvedIp string, err error) {
	bareHost := host
	port := "2719"
	if strings.Contains(host, ":") {
		parts := strings.SplitN(host, ":", 2)
		bareHost = parts[0]
		port = parts[1]
	}

	if net.ParseIP(bareHost) != nil {
		return bareHost + ":" + port, nil
	}

	dnsConn, err := net.Dial("tcp", "127.0.0.1:7007")
	if err != nil {
		return "", errors.New("Pengo DNS server is not responding")
	}
	dnsConn.Write([]byte(bareHost + "\n"))
	response, err := io.ReadAll(dnsConn)
	resolved := parseDnsResponse(string(response))
	return resolved, nil
}

func parseDnsResponse(response string) string {
	lines := strings.Split(response, "\n")
	status := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(status, "200") || len(lines) < 2 {
		return ""
	}
	return strings.TrimSpace(lines[1])
}

func main() {
	reader := bufio.NewReader(os.Stdin);
	fmt.Println("Enter pengo URI:");
	input, err := reader.ReadString('\n');
	if err != nil {
		log.Println(err);
	}

	host, request, err := parseInput(input);
	if err != nil {
		log.Fatalln(err)
	}
	resolvedIp, err := resolveDns(host)
	if err != nil {
		log.Println(err)
		return;
	}
	response := makeRequest(resolvedIp, request)
	fmt.Print(response)
}
