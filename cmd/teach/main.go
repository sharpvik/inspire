package main

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/sharpvik/env-go"
	"github.com/sharpvik/inspire/challenge"
	"github.com/sharpvik/inspire/handler"
	"github.com/sharpvik/inspire/message"
	"github.com/sharpvik/inspire/server"
)

var difficulty uint32

func init() {
	difficultyVar := env.GetOr("DIFFICULTY", "32")
	difficultyU64, err := strconv.ParseUint(difficultyVar, 10, 32)
	if err != nil {
		log.Fatalln(err)
	}
	difficulty = uint32(difficultyU64)
}

func handle() handler.Handler {
	return func(_ message.Message) message.Message {
		quote := quotes[rand.Intn(len(quotes))]
		log.Printf("quoting \"%s...\"", quote[:40])
		return message.Message(quote)
	}
}

func main() {
	challenge := challenge.WithDifficulty(difficulty)
	server := server.New(handle(), challenge)
	if err := server.ListenAndServe(":8910"); err != nil {
		log.Fatalln(err)
	}
}
