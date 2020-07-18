package main

import (
	"log"
	"net/url"
	"os"

	"github.com/tarstars/icfpc2020/diseaz/interpreter"
)

func main() {
	serverURL, err := url.Parse(os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	playerKey := os.Args[2]
	values := url.Values{}
	values.Add("apiKey", playerKey)
	serverURL.RawQuery = values.Encode()

	log.Printf("ServerUrl: %s", serverURL)

	c := interpreter.NewContext(serverURL)

	fn := os.Args[3]
	f, err := os.Open(fn)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	tok := interpreter.ParseReader(c, f)
	log.Printf("Result: %s", tok)
	log.Printf("Evals: %d", c.EvalCount)
}
