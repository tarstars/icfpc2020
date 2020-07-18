package interpreter

import (
	"fmt"
)

func Interpret(c Context, p Program) (Token, error) {
	pin := p.Clone()
	var pout Program
	for t, notEmpty := pin.Pops(); notEmpty; t, notEmpty = pin.Pops() {
		switch t.(type) {
		// case Var:
		// 	vp, err := tt.Value(c)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	pin.PushProgram(vp)
		case Ap:
			ft := pout.Pop()
			at := pout.Pop()
			pout.Push(Ap2{F: ft, A: at})
		default:
			pout.Push(t)
		}
	}
	if len(pout) != 1 {
		return nil, fmt.Errorf("Invalid compilation result: %s", pout)
	}
	return pout[0], nil
}
