package interpreter

import (
	"log"
)

type Context interface {
	GetVar(n int) Token
	SetVar(n int, v Token)
}

type Ctx struct {
	Vars map[int]Token
}

func NewContext() *Ctx {
	return &Ctx{
		Vars: make(map[int]Token),
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
