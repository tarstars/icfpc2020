package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/url"
	"os"
	"strconv"

	gx "github.com/tarstars/icfpc2020/diseaz/interpreter"
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
	Picture *gx.Picture `json:",inline"`
	Results []string    `json:""`
}

func (r *Result) AddResults(anss ...gx.Token) {
	for _, ans := range anss {
		r.Results = append(r.Results, ans.Galaxy())
	}
}

func main() {
	flag.Parse()

	serverURL, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Panic(err)
	}
	playerKey, err := strconv.ParseInt(flag.Arg(1), 2, 64)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("ServerUrl: %s; PlayerKey: %d", serverURL, playerKey)

	c := gx.NewContext(serverURL)
	r := &Result{}

	program := fmt.Sprintf("ap send (2, %s, nil)", playerKey)
	rs := gx.ParseString(c, program)
	r.AddResults(rs...)
	r.Picture = c.Picture()
	json.NewEncoder(os.Stdout).Encode(r)

	// if len(*drawOut) > 0 {
	// 	savePicture(*drawOut, c.Picture())
	// }
}
