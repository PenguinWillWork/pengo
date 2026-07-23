package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

func handleConnection(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	input, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	log.Println("Connection on 127.0.0.1:2719 from: " + connection.LocalAddr().String())
	log.Println(input)

	var body string
	var status string
	input = strings.TrimSpace(input)
	switch input {
	case "/home":
		status = "200 OK"
		body = `<div style="max-width:620px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;color:#1b2636;line-height:1.65;">
  <div style="font-size:.75rem;letter-spacing:.1em;text-transform:uppercase;color:#9aa2b1;margin-bottom:1.25rem;">pengo://home</div>
  <h1 style="font-size:1.6rem;font-weight:600;letter-spacing:-.01em;margin:0 0 1.1rem;">Pengo</h1>
  <p style="margin:0 0 1.1rem;">This is the first page displayed by pengo-browser, served over the Pengo protocol and resolved through Pengo DNS.</p>
  <p style="margin:0 0 1.6rem;color:#5b6472;">A small internet for penguins.</p>
  <div style="border-top:1px solid #e2e5ea;padding-top:1.1rem;">
    <a href="pengo://welcome" style="color:#1b2636;text-decoration:none;border-bottom:1px solid currentColor;">about &rarr;</a>
  </div>
</div>`
	default:
		status = "404 NOT FOUND"
		body = "unknown command: " + input + "\n"
	}
	response := makeResponse(status, body)
	connection.Write([]byte(response))
}

// standartizing response.. statuses need some structuring and the whole function requires some sanitizing before blindly appending things to the response i guess
func makeResponse(status string, body string) string {
	proto := "PENGO/0.1"
	contentLength := len(body)
	return proto + "\n" + status + "\n" +
		"Content-Length:" + strconv.Itoa(contentLength) + "\n\n" +
		body
}

func main() {
	listener, err := net.Listen("tcp", ":2719")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Listening on " + listener.Addr().String())
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}
