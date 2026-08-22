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

func FormRequest(method string, host string, headers []string, requestPath string) []byte {
	version := "PENGO/1.0";
	request := version + " " + method + " " + requestPath + " " + host + "\n" + strings.Join(headers, "\n") + "\n\n";
	return []byte(request)
}

func ParseRequest(reader *bufio.Reader) (Request, error) {
	request := readRequestInput(reader)
	log.Println("request: " + request)
	splitRequest := strings.Split(request, "\n")
	splitMainLine := strings.Split(splitRequest[0], " ")
	version := splitMainLine[0]
	method := splitMainLine[1]
	requestPath := splitMainLine[2]
	host := splitMainLine[3]
	if !strings.Contains(version,"PENGO") || method == "" || requestPath == "" || host == "" {
		return Request{}, errors.New("The request we got from the client is malformed")
	}
	headers := splitRequest[1:len(splitRequest)-2]
	return Request{version, method, requestPath, host, &headers}, nil
}

func readRequestInput(reader *bufio.Reader) string {
	var request strings.Builder
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading request line: " + err.Error())
		return ""
	}
	request.WriteString(line)

	newLineCount := 0
	for true {
		line, err = reader.ReadString('\n')
		if line == "\n" {
			newLineCount++
		} else {
			newLineCount = 0
		}
		if newLineCount == 2 {
			break
		}
		if err != nil {
			return ""
		}
		request.WriteString(line)
	}
	requestAsString := request.String()
	return requestAsString
}