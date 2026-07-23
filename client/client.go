package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type PengoResponse struct {
	Version string
	Status string
	ContentLength int
	Body  string
}

func parsePengoResponse(response string) (parsedResponse PengoResponse, err error) {
	//separating head from the body
	head, body, found := strings.Cut(response, "\n\n");
	
	//if not found -> error
	if !found {
		return PengoResponse{}, errors.New("The response that we got from the host is malformed: no header to body separator")
	}
	//chopping the headers
	headers := strings.Split(head, "\n");
	//if it ends up having less headers than needed (should at least have version, status, length lines) -> error
	if len(headers) < 3 {
		return PengoResponse{}, errors.New("The response that we got from the host is malformed: some headers are missing")
	} 
	//first two lines are by default version and status
	parsedResponse.Version = headers[0]
	parsedResponse.Status = headers[1];
	for i, h := range headers {
		//skip first 2 hearders - already asigned
		if i < 2 {
			continue;
		}
		header, value, found := strings.Cut(h, ":");
		if !found {
			return PengoResponse{}, errors.New("The response that we got from the host is malformed: some headers are missing")
		}
		switch header {
		case "Content-Length":
			length, err := strconv.Atoi(value)
			if err != nil {
				return PengoResponse{}, errors.New("malformed Content-Length: " + value)
			}
			parsedResponse.ContentLength = length
		}
	}
	if len(body) < parsedResponse.ContentLength {
		return PengoResponse{}, errors.New("truncated body: fewer bytes thanContent-Length promised")
	}
	parsedResponse.Body = body
	return parsedResponse, nil;
}

func makeRequest(ip string, request string) (PengoResponse, error) {
	log.Println("requesting: " + ip + " request: " + request + "\n")
	connection, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatalln("Couldn't connect to the pengo server on " + ip + "\n");
	}
	connection.Write([]byte(request + "\n"));
	response, err := io.ReadAll(connection);
	parsedResponse, err := parsePengoResponse(string(response));
	if err != nil {
		return PengoResponse{}, err;
	}
	return parsedResponse, nil;
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

//resolving dns by domain name input, if already and ip -> return ip. Otherwise dial dns server 
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
	//extract ip from the dns server response -> assign to resolved
	resolved := parseDnsResponse(string(response))
	return resolved, nil
}

func parseDnsResponse(response string) string {
	lines := strings.Split(response, "\n")
	status := strings.TrimSpace(lines[0])
	//checking if first line is 200
	if !strings.HasPrefix(status, "200") || len(lines) < 2 {
		return ""
	}
	//returning second line since we already checked it should be an IP
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
	response, err := makeRequest(resolvedIp, request)
	if err != nil {
		fmt.Print(err)
		return;
	}
	fmt.Print(response.Version + "\n" + response.Status + "\n" + response.Body + "\n")
}
