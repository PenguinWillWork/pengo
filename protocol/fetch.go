package protocol

import (
	"io"
	"net"
	"pengo/resolver"
)

func Fetch(method string, host string, reqPath string, headers []string) (Response, error) {
	response, err := makeRequest(method, host, reqPath, headers)
	if err != nil {
		return Response{}, err;
	}
	return response, nil;
}



func makeRequest(method string, host string, request string, headers []string) (Response, error) {
	resolvedIp, err := resolver.Resolve(host);
	if err != nil {
		return Response{}, err;
	}
	connection, err := net.Dial("tcp", resolvedIp)
	if err != nil {
		return Response{}, err;
	}
	formedRequest := FormRequest(method, host, headers, request)
	connection.Write(formedRequest);
	response, err := io.ReadAll(connection);
	parsedResponse, err := ParseResponse(string(response));
	if err != nil {
		return Response{}, err;
	}
	return parsedResponse, nil;
}
