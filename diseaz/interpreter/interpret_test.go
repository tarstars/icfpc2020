package interpreter

import (
	"log"
	"reflect"
	"testing"
)

func runProgram(t *testing.T, ts ...Token) Token {
	c := NewContext()
	p := NewProgram(ts...)
	tok, err := Interpret(c, p)
	if err != nil {
		t.Fatalf("Interpret failed: %s", err)
	}
	log.Printf("Program: %s", tok)
	r := tok.Eval(c)
	return r
}

func testProgram(t *testing.T, expected Token, ts ...Token) Token {
	r := runProgram(t, ts...)
	if !reflect.DeepEqual(r, expected) {
		t.Errorf("Result does not match: %s != %s", expected, r)
	}
	return r
}

func TestAdd(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Ap{}, Add{}, Int{V: 1}, Int{V: 41},
	)
}

func TestMul(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Ap{}, Mul{}, Int{V: 6}, Int{V: 7},
	)
}

func TestDiv(t *testing.T) {
	testProgram(t, Int{V: -6},
		Ap{}, Ap{}, Div{}, Int{V: 43}, Int{V: -7},
	)
}

func TestTrue(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Ap{}, True{}, Int{V: 42}, Int{V: 43},
	)
}

func TestFalse(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Ap{}, False{}, Int{V: 43}, Int{V: 42},
	)
}

func TestEq(t *testing.T) {
	testProgram(t, Int{V: 1},
		Ap{}, Ap{},

		// (
		Ap{}, Ap{}, Eq{},
		Int{V: 6},
		Ap{}, Ap{}, Div{}, Int{V: 42}, Int{V: 7},
		// )

		Int{V: 1},
		Int{V: 0},
	)
}

func TestPwr2(t *testing.T) {
	testProgram(t, Int{V: 1},
		Ap{}, Pwr2{}, Int{V: 0},
	)
	testProgram(t, Int{V: 2},
		Ap{}, Pwr2{}, Int{V: 1},
	)
	testProgram(t, Int{V: 16},
		Ap{}, Pwr2{}, Int{V: 4},
	)
}

func TestI(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Ap{}, I{}, Ap{}, Add{}, Int{V: 1}, Int{V: 41},
	)
}

func TestCar(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Car{},
		Ap{}, Ap{}, Cons{}, Int{V: 42}, Int{V: 41},
	)
}

func TestCdr(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Cdr{},
		Ap{}, Ap{}, Cons{}, Int{V: 41}, Int{V: 42},
	)
}

func TestIsNil(t *testing.T) {
	testProgram(t, True{},
		Ap{}, IsNil{}, Nil{},
	)
}

func TestIsNotNil(t *testing.T) {
	testProgram(t, False{},
		Ap{}, IsNil{},
		Ap{}, Ap{}, Cons{}, Int{V: 41}, Int{V: 42},
	)
}

func TestDraw(t *testing.T) {
	testProgram(t, Points{Point{X: 1, Y: 2}, Point{X: 41, Y: 42}},
		Ap{}, Draw{},
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 1}, Int{V: 2},
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 41}, Int{V: 42},
		Nil{},
	)
}

func TestMultipledraw(t *testing.T) {
	testProgram(t, Points{Point{X: 1, Y: 2}, Point{X: 41, Y: 42}, Point{X: 3, Y: 4}},
		Ap{}, Multipledraw{},
		Ap{}, Ap{}, Cons{},
		// (
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 1}, Int{V: 2},
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 41}, Int{V: 42},
		Nil{},
		// )
		Ap{}, Ap{}, Cons{},
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 3}, Int{V: 4},
		Nil{},
		Nil{},
	)
}
