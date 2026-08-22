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
	method, host, reqPath, headers, err := parseInput(input);
	if err != nil {
		return Response{}, err;
	}
	log.Println("method: " + method + " host: " + host + " request path: " + reqPath + " headers: " + strings.Join(headers, ","))
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

// pengo://joe/about  ->  host "joe", path "/about"
//   normalizedInput     pengo://joe/about
//   proto               pengo
//   inputWithoutProto   joe/about
//   uriBody             ["joe" "about"]
func parseInput(input string) (method string, host string, request string, headers []string, err error ) {
	normalizedInput := strings.TrimSpace(strings.ToLower(input));
	if !strings.Contains(normalizedInput, "headers:[") {
		return "", "", "", nil, errors.New("Invalid input format: missing headers")
	}
	splitInput := strings.Split(normalizedInput, " headers:");
	if len(splitInput) < 1 {
		return "", "", "", nil, errors.New("Invalid input format")
	}
	if len(splitInput) > 1 {
		splitInput[1] = strings.TrimSpace(strings.TrimPrefix(splitInput[1], "["));
		splitInput[1] = strings.TrimSpace(strings.TrimSuffix(splitInput[1], "]"));
		headers = strings.Split(splitInput[1], ",");
	}
		host = strings.Split(splitInput[0], " ")[1];
	splitInput[0] = strings.TrimSpace(splitInput[0]);
	method = strings.Split(splitInput[0], " ")[0];
	if method != "fetch" && method != "submit" {
		return "", "", "", nil, errors.New("Invalid method")
	}

	proto := strings.Split(host, "://")[0];
	hostWithoutProto := strings.Replace(host, proto + "://", "", -1);
	if (proto != "pengo") {
		err := "Requested URI is not a Pengo URI";
		return "", "", "", nil,errors.New(err);
	}

	uriBody := strings.SplitN(hostWithoutProto, "/", 2);
	hostFromUri := uriBody[0];
	path := "/";
	if len(uriBody) > 1 && uriBody[1] != "" {
		path = "/" + uriBody[1];
	}
	return method, hostFromUri, path, headers, nil;
}

//FETCH pengo://joe/about headers:[PebblePublicKey: abc123, PebbleTimestamp: 1234567890, PebbleNonce: 123456, PebbleSignature: def456] 
