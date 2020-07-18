package interpreter

import "fmt"

type Context interface {
	GetVar(n int) (Program, error)
	SetVar(n int, v Program)
}

type Ctx struct {
	Vars map[int]Program
}

func NewContext() *Ctx {
	return &Ctx{
		Vars: make(map[int]Program),
	}
}

func (c *Ctx) GetVar(n int) (Program, error) {
	p, exists := c.Vars[n]
	if !exists {
		return nil, fmt.Errorf("Variable does not exist: %d", n)
	}
	return p, nil
}

func (c *Ctx) SetVar(n int, p Program) {
	c.Vars[n] = p
}
