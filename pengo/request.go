package pengo

import (
	"bufio"
	"errors"
	"log"
	"strings"
)

type Request struct {
	Version string
	Method string
	RequestPath string
	Host string
	Headers *[]string
}

//PENGO/1.0 fetch / welcome
//pebblepublickey: abc123
//pebbletimestamp: 1234567890
//pebblenonce: 123456
//pebblesignature: def456
func FormRequest(method string, host string, headers []string, requestPath string) []byte {
	version := "PENGO/1.0"
	request := version + " " + method + " " + host + " "  + requestPath + "\n"
	if len(headers) > 0 {
		for _, h := range headers {
			strings.TrimSpace(h)
		}
		request += strings.Join(headers, "\n") + "\n"
	}
	request += "\n"
	return []byte(request)
}

func ParseRequest(reader *bufio.Reader) (Request, error) {
	parsedRequest := Request{}

	request := readRequestInput(reader)
	log.Println("request: " + request)
	splitRequest := strings.Split(request, "\n")
	splitMainLine := strings.Split(splitRequest[0], " ")

	parsedRequest.Version = splitMainLine[0]
	parsedRequest.RequestPath = splitMainLine[3]
	parsedRequest.Method = splitMainLine[1]
	parsedRequest.Host = splitMainLine[2]

	if !strings.Contains(parsedRequest.Version,"PENGO") || parsedRequest.Method == "" || parsedRequest.RequestPath == "" || parsedRequest.Host == "" {
		return Request{}, errors.New("The request we got from the client is malformed")
	}
	headers := splitRequest[1:len(splitRequest)-2]
	parsedRequest.Headers = &headers
	return parsedRequest, nil
}

func readRequestInput(reader *bufio.Reader) string {
	var request strings.Builder
	for {
		line, err := reader.ReadString('\n')


		if err != nil {
			return ""
		}
		request.WriteString(line)

		if strings.HasSuffix(request.String(), "\n\n") {
			break
		}
	}
	requestAsString := request.String()
	return requestAsString
}