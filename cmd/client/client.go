package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"pengo-proto/pengo"
)

func main() {
	reader := bufio.NewReader(os.Stdin);
	fmt.Fprint(os.Stderr, "Welcome to Pengo client!\n")
	fmt.Print("pengo> ");
	input, err := reader.ReadString('\n');
	if err != nil {
		log.Println(err);
	}

	response, err := pengo.Fetch(input)
	if err != nil {
		fmt.Print(err)
		return;
	}
	fmt.Print(response.Version + "\n" + response.Status + "\n" + string(response.Body) + "\n")
}
