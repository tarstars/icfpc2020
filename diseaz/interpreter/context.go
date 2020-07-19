package interpreter

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Context interface {
	GetVar(n int) Token
	SetVar(n int, v Token)
	Send(message string) string
}

type Ctx struct {
	Vars      map[int]Token
	ServerURL *url.URL
}

func NewContext(serverURL *url.URL) *Ctx {
	return &Ctx{
		Vars:      make(map[int]Token),
		ServerURL: serverURL,
	}
}

func (c *Ctx) GetVar(n int) Token {
	p, exists := c.Vars[n]
	if !exists {
		log.Panicf("Variable does not exist: %d", n)
	}
	return p
}

func (c *Ctx) SetVar(n int, p Token) {
	c.Vars[n] = p
}

func (c *Ctx) Send(message string) string {
	log.Printf("Send: %#v", message)
	res, err := http.Post(c.ServerURL.String(), "text/plain", strings.NewReader(message))
	if err != nil {
		log.Panicf("Unexpected server response:\n%v", err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panicf("Unexpected server response:\n%v", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Panicf("Unexpected server response:\nHTTP code: %d\nResponse body: %s", res.StatusCode, body)
	}

	r := string(body)
	log.Printf("Recv: %#v", r)
	return r
}
