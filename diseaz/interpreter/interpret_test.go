package interpreter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runProgram(t *testing.T, ts ...Token) Token {
	c := NewContext(nil)
	p := NewProgram(ts...)
	tok, err := Interpret(c, p)
	require.NoError(t, err, "Interpret failed")
	// log.Printf("Program: %s", tok)
	r := c.Eval(tok)
	return r
}

func testProgram(t *testing.T, expected Token, ts ...Token) Token {
	r := runProgram(t, ts...)
	assert.Equal(t, expected, r)
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

func TestParsePwr2(t *testing.T) {
	c := NewContext(nil)
	tok := ParseLine(c, ":42 = ap ap s ap ap c ap eq 0 1 ap ap b ap mul 2 ap ap b :42 ap add -1")
	assert.Nil(t, tok)
	tok = ParseLine(c, "ap :42 4")
	require.Len(t, tok, 1)
	assert.Equal(t, Int{V: 16}, tok[0])
}

func TestParsePwr2Reader(t *testing.T) {
	c := NewContext(nil)
	text := `:42 = ap ap s ap ap c ap eq 0 1 ap ap b ap mul 2 ap ap b :42 ap add -1
ap :42 4`
	tok := ParseReader(c, strings.NewReader(text))
	require.Len(t, tok, 1)
	assert.Equal(t, Int{V: 16}, tok[0])
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
	testProgram(t, NewPicture(Point{X: 1, Y: 2}, Point{X: 41, Y: 42}),
		Ap{}, Draw{},
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 1}, Int{V: 2},
		Ap{}, Ap{}, Cons{}, Ap{}, Ap{}, Vec{}, Int{V: 41}, Int{V: 42},
		Nil{},
	)
}

func TestDrawListCons(t *testing.T) {
	c := NewContext(nil)
	text := "ap draw (ap ap vec 1 2, ap ap vec 41 42)"
	tok := ParseReader(c, strings.NewReader(text))
	require.Len(t, tok, 1)
	assert.Equal(t, NewPicture(Point{X: 1, Y: 2}, Point{X: 41, Y: 42}), tok[0])
}

func TestMultipledraw(t *testing.T) {
	testProgram(t, NewPicture(Point{X: 1, Y: 2}, Point{X: 41, Y: 42}).DrawPts(Point{X: 3, Y: 4}),
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

func TestIf0(t *testing.T) {
	testProgram(t, Int{V: 42},
		Ap{}, Ap{}, Ap{}, If0{},
		Ap{}, Dec{}, Int{V: 1},
		Ap{}, Ap{}, Add{}, Int{V: 1}, Int{V: 41},
		Nil{},
	)
}

func TestModulateInt(t *testing.T) {
	assert.Equal(t, "010", modInt(0))
	assert.Equal(t, "01100001", modInt(1))
	assert.Equal(t, "10100001", modInt(-1))
	assert.Equal(t, "0111000101010", modInt(42))
}

func TestModulate(t *testing.T) {
	assert.Equal(t, "00", ModulateToken(Nil{}))
	assert.Equal(t, "110000", ModulateToken(Cons2{X0: Nil{}, X1: Nil{}}))
	assert.Equal(t, "110110000101100010", ModulateToken(Cons2{X0: Int{V: 1}, X1: Int{V: 2}}))
}

func TestDemodulate(t *testing.T) {
	v := DemodulateToken("00")
	assert.Equal(t, "00", ModulateToken(v))

	v = DemodulateToken("110000")
	assert.Equal(t, "110000", ModulateToken(v))

	v = DemodulateToken("110110000101100010")
	assert.Equal(t, "110110000101100010", ModulateToken(v))
}

func TestModDemod(t *testing.T) {
	testProgram(t, Signal{S: "00"},
		Ap{}, Modulate{},
		Ap{}, Demodulate{}, Signal{S: "00"},
	)
	testProgram(t, Signal{S: "110000"},
		Ap{}, Modulate{},
		Ap{}, Demodulate{}, Signal{S: "110000"},
	)
	testProgram(t, Signal{S: "110110000101100010"},
		Ap{}, Modulate{},
		Ap{}, Demodulate{}, Signal{S: "110110000101100010"},
	)
}
