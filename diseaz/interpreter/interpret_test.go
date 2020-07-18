package interpreter

import (
	"reflect"
	"testing"
)

func runProgram(t *testing.T, ts ...Token) Program {
	c := NewContext()
	pout, err := Interpret(c, NewProgram(ts...))
	if err != nil {
		t.Fatal(err)
	}
	return pout
}

func testProgram(t *testing.T, expected Program, ts ...Token) Program {
	c := NewContext()
	pout, err := Interpret(c, NewProgram(ts...))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(pout, expected) {
		t.Fail()
	}
	return pout
}

func TestAdd(t *testing.T) {
	testProgram(t, Program{Int{V: 42}},
		Ap{},
		Ap{},
		Add{},
		Int{V: 1},
		Int{V: 41},
	)
}

func TestMul(t *testing.T) {
	testProgram(t, Program{Int{V: 42}},
		Ap{},
		Ap{},
		Mul{},
		Int{V: 6},
		Int{V: 7},
	)
}

func TestEq(t *testing.T) {
	testProgram(t, Program{True{}},
		Ap{},
		Ap{},
		Eq{},
		Int{V: 6},
		Ap{},
		Ap{},
		Div{},
		Int{V: 42},
		Int{V: 7},
	)
}
