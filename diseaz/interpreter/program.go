package interpreter

import (
	"fmt"
	"strings"
)

type Program []Token

func NewProgram(ts ...Token) Program {
	return ts
}

func (p Program) String() string {
	var r []string
	for _, t := range p {
		switch tt := t.(type) {
		case fmt.Stringer:
			r = append(r, tt.String())
		default:
			r = append(r, fmt.Sprintf("%#v", t))
		}
	}
	return "[" + strings.Join(r, " ") + "]"
}

func (p Program) Clone() Program {
	r := make(Program, len(p))
	copy(r, p)
	return r
}

func (p Program) Reverse() Program {
	pout := make(Program, len(p))
	for i := range p {
		pout[len(pout)-i-1] = p[i]
	}
	return pout
}

func (p *Program) Pops() (Token, bool) {
	if len(*p) == 0 {
		return nil, false
	}
	n := *p
	*p = (*p)[:len(n)-1]
	return n[len(n)-1], true
}

func (p *Program) Pop() Token {
	r, ok := p.Pops()
	if !ok {
		panic(fmt.Errorf("Pop from empty program"))
	}
	return r
}

func (p *Program) Push(t Token) {
	*p = append(*p, t)
}

func (p *Program) PushProgram(ts Program) {
	for _, t := range ts {
		p.Push(t)
	}
}
