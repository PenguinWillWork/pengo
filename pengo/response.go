package pengo

import (
	"errors"
	"strconv"
	"strings"
)

type Response struct {
	Version string
	Status string
	ContentLength int
	Body  string
}

const (
	StatusOK = "200 OKY"
	StatusNotFound = "400 CONTENT NOT FOUND"
	StatusInternalServerError = "500 SERVER DID SOMETHING WRONG"
)

// standartizing response.. statuses need some structuring and the whole function requires some sanitizing before blindly appending things to the response i guess
func MakeResponse(status string, body string) string {
	proto := "PENGO/0.1"
	contentLength := len(body)
	return proto + "\n" + status + "\n" +
		"Content-Length:" + strconv.Itoa(contentLength) + "\n\n" +
		body
}

func ParseResponse(response string) (parsedResponse Response, err error) {
	//separating head from the body
	head, body, found := strings.Cut(response, "\n\n");

	//if not found -> error
	if !found {
		return Response{}, errors.New("The response that we got from the host is malformed: no header to body separator")
	}
	//chopping the headers
	headers := strings.Split(head, "\n");
	//if it ends up having less headers than needed (should at least have version, status, length lines) -> error
	if len(headers) < 3 {
		return Response{}, errors.New("The response that we got from the host is malformed: some headers are missing")
	}
	//first two lines are by default version and status
	parsedResponse.Version = headers[0];
	parsedResponse.Status = headers[1];
	for i, h := range headers {
		//skip first 2 hearders - already asigned
		if i < 2 {
			continue;
		}
		header, value, found := strings.Cut(h, ":");
		if !found {
			return Response{}, errors.New("The response that we got from the host is malformed: some headers are missing")
		}
		switch header {
		case "Content-Length":
			length, err := strconv.Atoi(value)
			if err != nil {
				return Response{}, errors.New("malformed Content-Length: " + value)
			}
			parsedResponse.ContentLength = length
		}
	}
	if len(body) < parsedResponse.ContentLength {
		return Response{}, errors.New("truncated body: fewer bytes thanContent-Length promised")
	}
	parsedResponse.Body = body
	return parsedResponse, nil;
}
