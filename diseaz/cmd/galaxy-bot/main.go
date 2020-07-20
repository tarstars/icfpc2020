package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	gx "github.com/tarstars/icfpc2020/diseaz/interpreter"
)

type ICons = gx.ICons
type Int = gx.Int

type Result struct {
	Picture *gx.Picture `json:",inline"`
	Results []string    `json:""`
}

func (r *Result) AddResults(anss ...gx.Token) {
	for _, ans := range anss {
		r.Results = append(r.Results, ans.Galaxy())
	}
}

func command(c gx.Context, name string, program string) gx.Token {
	logr := &Result{}
	rs := gx.ParseString(c, program)
	logr.AddResults(rs...)
	logr.Picture = c.Picture()
	json.NewEncoder(os.Stdout).Encode(logr)

	status := rs[0].(ICons).Car().(Int).V
	if status == 0 {
		log.Panic(fmt.Errorf("%s failed", name))
	}

	return rs[0]
}

type GameStage int

const (
	GamePending  GameStage = 0
	GameStarted  GameStage = 1
	GameFinished GameStage = 2
)

func main() {
	flag.Parse()

	serverURL, err := url.Parse(flag.Arg(0))
	if err != nil {
		log.Panic(err)
	}
	serverURL.Path = "/aliens/send"
	playerKey, err := strconv.ParseInt(flag.Arg(1), 10, 64)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("ServerUrl: %s; PlayerKey: %d", serverURL, playerKey)

	c := gx.NewContext(serverURL)
	joinResponse := command(c, "JOIN", fmt.Sprintf("ap send (2, %d, nil)", playerKey)).(ICons)
	joinResponse1 := joinResponse.Cdr().(ICons)
	stage := GameStage(joinResponse1.Car().(Int).V)
	log.Printf("GameStage: %s", stage)

	command(c, "START", fmt.Sprintf("ap send (3, %d, (1, 1, 1, 1))", playerKey))
}
