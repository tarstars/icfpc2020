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

	Eval(t Token) Token

	Send(message string) string
	SendToken(v Token) Token

	Picture() *Picture
	// OutFormat() OutFormat

	// Level() int
	// Enter() int
	// Leave() int
	// CountEval() int
}

type OutFormat string

const (
	OutFormatS OutFormat = "s"
	OutFormatG OutFormat = "g"
)

type Ctx struct {
	Vars      map[int]Token
	ServerURL *url.URL
	Pic       *Picture
	CallLevel int
	EvalCount int
}

func NewContext(serverURL *url.URL) *Ctx {
	return &Ctx{
		Vars:      make(map[int]Token),
		ServerURL: serverURL,
		Pic:       NewPicture(),
	}
}

func (c Ctx) GetVar(n int) Token {
	p, exists := c.Vars[n]
	if !exists {
		log.Panicf("Variable does not exist: %d", n)
	}
	return p
}

func (c *Ctx) SetVar(n int, p Token) {
	c.Vars[n] = p
}

func (c *Ctx) SendToken(v Token) Token {
	return DemodulateToken(c.Send(ModulateToken(v)))
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

func (c Ctx) Level() int {
	return c.CallLevel
}

func (c *Ctx) Enter() int {
	c.CallLevel++
	return c.CallLevel
}

func (c *Ctx) Leave() int {
	c.CallLevel--
	return c.CallLevel
}

func (c *Ctx) CountEval() int {
	c.EvalCount++
	return c.EvalCount
}

func (c *Ctx) Picture() *Picture {
	return c.Pic
}

func (c Ctx) OutFormat() OutFormat {
	return OutFormatG
}
