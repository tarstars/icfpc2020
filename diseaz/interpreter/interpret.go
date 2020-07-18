package interpreter

import (
	"fmt"
)

func Interpret(c Context, p Program) (Program, error) {
	pin := p.Clone()
	var pout Program
	for t, notEmpty := pin.Pops(); notEmpty; t, notEmpty = pin.Pops() {
		switch tt := t.(type) {
		case Var:
			vp, err := c.GetVar(tt.N)
			if err != nil {
				return nil, err
			}
			pin.PushProgram(vp)
		case Ap:
			t1 := pout.Pop()
			f, ok := t1.(Func)
			if !ok {
				return nil, fmt.Errorf("Apply non-function %s", t1)
			}
			t2 := pout.Pop()
			ap, err := f.Apply(t2)
			if err != nil {
				return nil, err
			}
			pin.PushProgram(ap)
		default:
			pout.Push(t)
		}
	}
	return pout, nil
}
