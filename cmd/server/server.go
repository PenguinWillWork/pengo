package main

import (
	"bufio"
	"log"
	"net"
	"pengo-proto/pengo"
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
		status = pengo.StatusOK
		body = `<div style="max-width:620px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;color:#1b2636;line-height:1.65;">
  <div style="font-size:.75rem;letter-spacing:.1em;text-transform:uppercase;color:#9aa2b1;margin-bottom:1.25rem;">pengo://welcome</div>
  <h1 style="font-size:1.6rem;font-weight:600;letter-spacing:-.01em;margin:0 0 1.1rem;">Pengo</h1>
  <p style="margin:0 0 1.1rem;">This is the first page displayed by pengo-browser, served over the Pengo protocol and resolved through Pengo DNS.</p>
  <p style="margin:0 0 1.6rem;color:#5b6472;">A small internet for penguins.</p>
  <div style="border-top:1px solid #e2e5ea;padding-top:1.1rem;">
    <a href="pengo://welcome" style="color:#1b2636;text-decoration:none;border-bottom:1px solid currentColor;">about &rarr;</a>
  </div>
</div>`
	default:
		status = pengo.StatusNotFound
		body = `<div style="max-width:620px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;color:#1b2636;line-height:1.65;">
  <div style="font-size:.75rem;letter-spacing:.1em;text-transform:uppercase;color:#9aa2b1;margin-bottom:1.25rem;">pengo://404</div>
  <h1 style="font-size:1.6rem;font-weight:600;letter-spacing:-.01em;margin:0 0 1.1rem;">Not found</h1>
  <p style="margin:0 0 1.1rem;">The page <strong style="font-weight:600;">` + input + `</strong> does not exist on this Pengo server.</p>
  <p style="margin:0 0 1.6rem;color:#5b6472;">The penguins looked everywhere. It is not here.</p>
  <div style="border-top:1px solid #e2e5ea;padding-top:1.1rem;">
    <a href="pengo://home" style="color:#1b2636;text-decoration:none;border-bottom:1px solid currentColor;">&larr; home</a>
  </div>
</div>`
	}
	response := pengo.MakeResponse(status, body)
	connection.Write([]byte(response))
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
