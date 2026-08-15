package resolver

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func Resolve(host string) (resolvedIp string, err error) {
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
		return "", fmt.Errorf("connect to pengo dns server: %w", err)
	}
	defer dnsConn.Close();

	dnsConn.Write([]byte(bareHost + "\n"))
	response, err := io.ReadAll(dnsConn)
	if err != nil {
		return "", fmt.Errorf("read from pengo dns server: %w", err)
	}
	//extract ip from the dns server response -> assign to resolved
	resolved, err := parseDnsResponse(string(response))
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", host, err)
	}
	return resolved, nil
}

func parseDnsResponse(response string) (string, error) {
	lines := strings.Split(response, "\n")
	status := strings.TrimSpace(lines[0])
	//checking if first line is 200
	if !strings.HasPrefix(status, "200") || len(lines) < 2 {
		return "", errors.New("dns returned: " + status)
	}
	//returning second line since we already checked it should be an IP
	return strings.TrimSpace(lines[1]), nil
}
