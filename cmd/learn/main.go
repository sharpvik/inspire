package main

import (
	"fmt"
	"log"

	"github.com/sharpvik/env-go"
	"github.com/sharpvik/inspire/client"
	"github.com/sharpvik/inspire/message"
)

var (
	hey  = message.Message([]byte("hey"))
	host = env.GetOr("HOST", "localhost")
	addr = host + ":8910"
)

func main() {
	for client := client.New(addr); ; {
		response, err := client.Send(hey)
		abortOrError(err)
		fmt.Println(string(response))
	}
}

func abortOrError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
