package resolver

import (
	"errors"
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
		return "", errors.New("Pengo DNS server is not responding")
	}
	dnsConn.Write([]byte(bareHost + "\n"))
	response, err := io.ReadAll(dnsConn)
	//extract ip from the dns server response -> assign to resolved
	resolved := parseDnsResponse(string(response))
	return resolved, nil
}

func parseDnsResponse(response string) string {
	lines := strings.Split(response, "\n")
	status := strings.TrimSpace(lines[0])
	//checking if first line is 200
	if !strings.HasPrefix(status, "200") || len(lines) < 2 {
		return ""
	}
	//returning second line since we already checked it should be an IP
	return strings.TrimSpace(lines[1])
}
