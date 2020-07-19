package main

import (
	"encoding/json"
	"flag"
	"image"
	"image/png"
	"log"
	"net/url"
	"os"

	"github.com/tarstars/icfpc2020/diseaz/interpreter"
)

func savePicture(fn string, pic image.Image) {
	f, err := os.Create(fn)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	err = png.Encode(f, pic)
	if err != nil {
		log.Panic(err)
	}
}

type Result struct {
	Picture *interpreter.Picture `json:",inline"`
	Results []string             `json:""`
}

func main() {
	server := flag.String("server", "https://icfpc2020-api.testkontur.ru/aliens/send", "Server URL")
	key := flag.String("key", "faa0647bb89f42d6a0a1850cf1b71954", "Player key")
	drawOut := flag.String("draw", "", "Output picture file")
	flag.Parse()

	serverURL, err := url.Parse(*server)
	if err != nil {
		log.Panic(err)
	}
	values := url.Values{}
	values.Add("apiKey", *key)
	serverURL.RawQuery = values.Encode()
	log.Printf("ServerUrl: %s", serverURL)

	c := interpreter.NewContext(serverURL)

	r := Result{}
	for _, fn := range flag.Args() {
		f, err := os.Open(fn)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		toks := interpreter.ParseReader(c, f)
		for _, tok := range toks {
			r.Results = append(r.Results, tok.Galaxy())
			log.Printf("Result(s): %s", tok)
			log.Printf("Result: %s", tok.Galaxy())
		}
	}
	log.Printf("Evals: %d", c.EvalCount)

	r.Picture = c.Picture()
	json.NewEncoder(os.Stdout).Encode(r)

	if len(*drawOut) > 0 {
		savePicture(*drawOut, c.Picture())
	}
}
