package pengo

import (
	"bufio"
	"errors"
	"strings"
)

type Request struct {
	Version string
	RequestPath string
	Host string
}

func FormRequest(host string, requestPath string) []byte {
	version := "PENGO/1.0";
	request := version + "\n\n" + requestPath + "\n" + "Host:" + host + "\n"
	return []byte(request)
}

func ParseRequest(reader *bufio.Reader) (Request, error) {
	request := readRequestInput(reader)
	version, data, found := strings.Cut(request, "\n\n");
	requestData := strings.Split(data, "\n")
	if !found || !strings.Contains(version,"PENGO") || len(requestData) != 2 {
		return Request{}, errors.New("The request we got from the client is malformed")
	}
	if (requestData[1] == "") {
		requestData[1] = "/"	
	}
	return Request{version, requestData[0], requestData[1]}, nil
}

func readRequestInput(reader *bufio.Reader) string {
	var request strings.Builder
	for i := 0; i < 3; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
				return ""
		}
		request.WriteString(line)
	}
	requestAsString := request.String()
	return requestAsString 
}