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
	host, reqPath, err := parseUserInput(input);
	if err != nil {
		return Response{}, err;
	}
	response, err := makeRequest(host, reqPath)
	if err != nil {
		return Response{}, err;
	}
	return response, nil;
}



func makeRequest(host string, request string) (Response, error) {
	resolvedIp, err := resolver.Resolve(host);
	if err != nil {
		return Response{}, err;
	}
	log.Println("requesting: " + resolvedIp + " request: " + request + "\n")
	connection, err := net.Dial("tcp", resolvedIp)
	if err != nil {
		return Response{}, err;
	}
	formedRequest := FormRequest(host, request)
	connection.Write(formedRequest);
	response, err := io.ReadAll(connection);
	parsedResponse, err := ParseResponse(string(response));
	if err != nil {
		return Response{}, err;
	}
	return parsedResponse, nil;
}

func parseUserInput(input string) (host string, request string, err error){
	normalizedInput := strings.TrimSpace(strings.ToLower(input));
	proto := strings.Split(normalizedInput, "://")[0];
	inputWithoutProto := strings.Replace(normalizedInput, proto + "://", "", -1);
	if (proto != "pengo") {
		err := "Requested URI is not a Pengo URI";
		return "", "", errors.New(err);
	}

	uriBody := strings.SplitN(inputWithoutProto, "/", 2);
	hostFromUri := uriBody[0];
	path := "/";
	if len(uriBody) > 1 && uriBody[1] != "" {
		path = "/" + uriBody[1];
	}
	return hostFromUri, path, nil;
}
