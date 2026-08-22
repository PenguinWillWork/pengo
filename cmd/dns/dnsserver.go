package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strings"
)

func resolveDnsRequest(connection net.Conn, registry map[string]string) {
	defer connection.Close()

	requester := connection.RemoteAddr().String()

	reader := bufio.NewReader(connection)
	input, err := reader.ReadString('\n')
	log.Println(requester + " sent request: " + input)
	if err != nil {
		log.Println(requester + " sent an unreadable request: " + err.Error())
		return
	}

	//input just has "domain name"
	domain := strings.TrimSpace(strings.ToLower(input))
	log.Println(requester + " is resolving " + domain)

	//look-up dns registry for domain passed by user
	ip, found := registry[domain]
	var response string
	if found {
		response = "200 OK\n" + ip + "\n"
		log.Println(requester + " resolved " + domain + " -> " + ip)
	} else {
		response = "404 NOT FOUND\n"
		log.Println(requester + " failed to resolve " + domain)
	}
	connection.Write([]byte(response))
}

func main() {
	//read registry file
	file, err := os.Open("registry.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var registry map[string]string
	//decode registry json
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&registry)
	if err != nil {
		log.Fatalln("could not read registry.json: " + err.Error())
	}
	log.Printf("loaded %d records from registry.json", len(registry))

	//create dns server listener on 7007
	listener, err := net.Listen("tcp", ":7007")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	log.Println("DNS server is listening at " + listener.Addr().String())

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go resolveDnsRequest(connection, registry)
	}
}
