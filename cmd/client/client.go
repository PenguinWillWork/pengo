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
	fmt.Println("Enter pengo URI:");
	input, err := reader.ReadString('\n');
	if err != nil {
		log.Println(err);
	}

	response, err := pengo.Fetch(input)
	if err != nil {
		fmt.Print(err)
		return;
	}
	fmt.Print(response.Version + "\n" + response.Status + "\n" + response.Body + "\n")
}
