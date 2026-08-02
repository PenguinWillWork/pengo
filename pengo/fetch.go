package pengo

import (
	"errors"
	"io"
	"log"
	"net"
	"pengo-proto/resolver"
	"strings"
)

func Fetch(input string) (Response, error) {
	host, request, err := parseInput(input);
	if err != nil {
		return Response{}, err;
	}
	resolvedIp, err := resolver.Resolve(host);
	if err != nil {
		return Response{}, err;
	}
	response, err := makeRequest(resolvedIp, request)
	if err != nil {
		return Response{}, err;
	}
	return response, nil;
}

func makeRequest(ip string, request string) (Response, error) {
	log.Println("requesting: " + ip + " request: " + request + "\n")
	connection, err := net.Dial("tcp", ip)
	if err != nil {
		return Response{}, errors.New("Couldn't connect to the pengo server on " + ip + "\n");
	}
	connection.Write([]byte(request + "\n"));
	response, err := io.ReadAll(connection);
	parsedResponse, err := ParseResponse(string(response));
	if err != nil {
		return Response{}, err;
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
	path := "/";
	if len(uriBody) > 1 && uriBody[1] != "" {
		path = "/" + uriBody[1];
	}
	return host, path, nil;
}
