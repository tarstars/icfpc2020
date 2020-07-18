package interpreter

import (
	"fmt"
)

func Interpret(c Context, p Program) (Token, error) {
	pin := p.Clone()
	var pout Program
	for t, notEmpty := pin.Pops(); notEmpty; t, notEmpty = pin.Pops() {
		// log.Printf("Token: %s", t)
		switch t.(type) {
		case Ap:
			ft := pout.Pop()
			at := pout.Pop()
			ff, ok := ft.(Func)
			if ok {
				pout.Push(ff.Apply(at))
			} else {
				pout.Push(Ap2{F: ft, A: at})
			}
		default:
			pout.Push(t)
		}
		// log.Printf("Out: %s", pout)
	}
	if len(pout) != 1 {
		return nil, fmt.Errorf("Invalid compilation result: %s", pout)
	}
	return pout[0], nil
}
