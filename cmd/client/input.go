package main

import (
	"errors"
	"strings"
)

// pengo://joe/about  ->  host "joe", path "/about"
//FETCH pengo://joe/about headers:[PebblePublicKey: abc123, PebbleTimestamp: 1234567890, PebbleNonce: 123456, PebbleSignature: def456] 
func parseInput(input string) (method string, host string, request string, headers []string, err error ) {
	normalizedInput := strings.TrimSpace(strings.ToLower(input));
	if !strings.Contains(normalizedInput, "headers:[") {
		return "", "", "", nil, errors.New("Invalid input format: missing headers\n")
	}
	splitInput := strings.Split(normalizedInput, " headers:");
	if len(splitInput) < 1 {
		return "", "", "", nil, errors.New("Invalid input format\n")
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
		return "", "", "", nil, errors.New("Invalid method\n")
	}

	proto := strings.Split(host, "://")[0];
	hostWithoutProto := strings.Replace(host, proto + "://", "", -1);
	if (proto != "pengo") {
		err := "Requested URI is not a Pengo URI\n";
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
