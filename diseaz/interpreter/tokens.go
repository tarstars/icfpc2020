package interpreter

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Token interface {
	Eval(c Context) Token
	String() string
}

type Value interface {
	Token
	Value() int64
}

type Var interface {
	Token
	Get(c Context) Token
}

type Func interface {
	Token
	Apply(v Token) Token
}

type Ap struct{}

func (t Ap) Eval(c Context) Token {
	return t
}

func (t Ap) String() string {
	return "ap"
}

type Ap2 struct {
	F Token
	A Token
}

func (t Ap2) Apply(v Token) Token {
	return Ap2{F: t, A: v}
}

func (t Ap2) Eval(c Context) Token {
	f := t.F.Eval(c).(Func)
	r := f.Apply(t.A).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Ap2) String() string {
	return fmt.Sprintf("(ap %s %s)", t.F, t.A)
}

type VarN struct {
	N int
}

func (t VarN) Eval(c Context) Token {
	return c.GetVar(t.N)
}

func (t VarN) String() string {
	return fmt.Sprintf(":%d", t.N)
}

type Int struct {
	V int64
}

func (t Int) Value() int64 {
	return t.V
}

func (t Int) Eval(c Context) Token {
	return t
}

func (t Int) String() string {
	return strconv.FormatInt(t.V, 10)
}

type Inc struct{}
type Inc1 struct {
	X0 Token
}

func (t Inc) Apply(v Token) Token {
	return Inc1{X0: v}
}

func (t Inc1) Eval(c Context) Token {
	v := t.X0.Eval(c).(Value)
	r := Int{V: v.Value() + 1}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Inc) Eval(c Context) Token {
	return t
}

func (t Inc) String() string {
	return "inc"
}

func (t Inc1) String() string {
	return fmt.Sprintf("(inc %s)", t.X0)
}

type Dec struct{}
type Dec1 struct {
	X0 Token
}

func (t Dec) Apply(v Token) Token {
	return Dec1{X0: v}
}

func (t Dec1) Eval(c Context) Token {
	v := t.X0.Eval(c).(Value)
	r := Int{V: v.Value() - 1}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Dec) Eval(c Context) Token {
	return t
}

func (t Dec) String() string {
	return "dec"
}

func (t Dec1) String() string {
	return fmt.Sprintf("(dec %s)", t.X0)
}

type Add struct{}
type Add1 struct {
	X0 Token
}
type Add2 struct {
	X0 Token
	X1 Token
}

func (t Add) Apply(v Token) Token {
	return Add1{X0: v}
}

func (t Add1) Apply(v Token) Token {
	return Add2{X0: t.X0, X1: v}
}

func (t Add2) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Value)
	x1 := t.X1.Eval(c).(Value)
	r := Int{V: x0.Value() + x1.Value()}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Add) Eval(c Context) Token {
	return t
}

func (t Add1) Eval(c Context) Token {
	return t
}

func (t Add) String() string {
	return "add"
}

func (t Add1) String() string {
	return fmt.Sprintf("(add %s)", t.X0)
}

func (t Add2) String() string {
	return fmt.Sprintf("(add %s %s)", t.X0, t.X1)
}

type Mul struct{}
type Mul1 struct {
	X0 Token
}
type Mul2 struct {
	X0 Token
	X1 Token
}

func (t Mul) Apply(v Token) Token {
	return Mul1{X0: v}
}

func (t Mul1) Apply(v Token) Token {
	return Mul2{X0: t.X0, X1: v}
}

func (t Mul2) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Value)
	x1 := t.X1.Eval(c).(Value)
	r := Int{V: x0.Value() * x1.Value()}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Mul) Eval(c Context) Token {
	return t
}

func (t Mul1) Eval(c Context) Token {
	return t
}

func (t Mul) String() string {
	return "mul"
}

func (t Mul1) String() string {
	return fmt.Sprintf("(mul %s)", t.X0)
}

func (t Mul2) String() string {
	return fmt.Sprintf("(mul %s %s)", t.X0, t.X1)
}

type Div struct{}
type Div1 struct {
	X0 Token
}
type Div2 struct {
	X0 Token
	X1 Token
}

func (t Div) Apply(v Token) Token {
	return Div1{X0: v}
}

func (t Div1) Apply(v Token) Token {
	return Div2{X0: t.X0, X1: v}
}

func (t Div2) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Value)
	x1 := t.X1.Eval(c).(Value)
	r := Int{V: x0.Value() / x1.Value()}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Div) Eval(c Context) Token {
	return t
}

func (t Div1) Eval(c Context) Token {
	return t
}

func (t Div) String() string {
	return "div"
}

func (t Div1) String() string {
	return fmt.Sprintf("(div %s)", t.X0)
}

func (t Div2) String() string {
	return fmt.Sprintf("(div %s %s)", t.X0, t.X1)
}

type Eq struct{}
type Eq1 struct {
	X0 Token
}
type Eq2 struct {
	X0 Token
	X1 Token
}

func (t Eq) Apply(v Token) Token {
	return Eq1{X0: v}
}

func (t Eq1) Apply(v Token) Token {
	return Eq2{X0: t.X0, X1: v}
}

func (t Eq2) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Value)
	x1 := t.X1.Eval(c).(Value)
	var r Token
	if x0.Value() == x1.Value() {
		r = True{}
	} else {
		r = False{}
	}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Eq) Eval(c Context) Token {
	return t
}

func (t Eq1) Eval(c Context) Token {
	return t
}

func (t Eq) String() string {
	return "eq"
}

func (t Eq1) String() string {
	return fmt.Sprintf("(eq %s)", t.X0)
}

func (t Eq2) String() string {
	return fmt.Sprintf("(eq %s %s)", t.X0, t.X1)
}

type Lt struct{}
type Lt1 struct {
	X0 Token
}
type Lt2 struct {
	X0 Token
	X1 Token
}

func (t Lt) Apply(v Token) Token {
	return Lt1{X0: v}
}

func (t Lt1) Apply(v Token) Token {
	return Lt2{X0: t.X0, X1: v}
}

func (t Lt2) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Value)
	x1 := t.X1.Eval(c).(Value)
	var r Token = False{}
	if x0.Value() < x1.Value() {
		r = True{}
	}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Lt) Eval(c Context) Token {
	return t
}

func (t Lt1) Eval(c Context) Token {
	return t
}

func (t Lt) String() string {
	return "lt"
}

func (t Lt1) String() string {
	return fmt.Sprintf("(lt %s)", t.X0)
}

func (t Lt2) String() string {
	return fmt.Sprintf("(lt %s %s)", t.X0, t.X1)
}

func mod(c Context, v Token) string {
	switch tt := v.Eval(c).(type) {
	case Int:
		return modInt(tt.V)
	case ICons:
		return modCons(c, tt)
	default:
		log.Panicf("Invalid `modulate` argument: %s", v)
	}
	return "ERROR"
}

func modInt(v int64) string {
	if v == 0 {
		return "010"
	}
	prefix := "01"
	if v < 0 {
		prefix = "10"
		v = -v
	}
	rs := []string{prefix}
	s := strconv.FormatInt(v, 2)
	n := (len(s) + 3) / 4
	rs = append(rs, strings.Repeat("1", n), "0")
	m := (4 * n) - len(s)
	if m > 0 {
		rs = append(rs, strings.Repeat("0", m))
	}
	rs = append(rs, s)
	return strings.Join(rs, "")
}

func modCons(c Context, v ICons) string {
	if v.IsNil() {
		return "00"
	}
	return "11" + mod(c, v.Car()) + mod(c, v.Cdr())
}

func demod(v string) (Token, string) {
	if len(v) == 0 {
		return nil, v
	}
	prefix, w := v[0:2], v[2:]
	if prefix == "00" {
		return Nil{}, w
	}
	if prefix == "11" {
		car, w := demod(w)
		cdr, w := demod(w)
		return Cons2{X0: car, X1: cdr}, w
	}
	return demodInt(v)
}

func demodInt(v string) (Token, string) {
	var negative bool
	prefix, w := v[0:2], v[2:]
	switch prefix {
	case "01":
		negative = false
	case "10":
		negative = true
	default:
		log.Panicf("Invalid modulated int prefix: %#v", v)
	}
	nlen := 0
	for ; w[0] == '1'; w = w[1:] {
		nlen += 4
	}
	w = w[1:]
	if nlen == 0 {
		return Int{V: 0}, w
	}
	num, w := w[:nlen], w[nlen:]
	n, err := strconv.ParseInt(num, 2, 64)
	if err != nil {
		log.Panic(err)
	}
	if negative {
		n = -n
	}
	return Int{V: n}, w
}

type Signal struct {
	S string
}

func (t Signal) Eval(c Context) Token {
	return t
}

func (t Signal) String() string {
	return fmt.Sprintf("%#v", t.S)
}

type Modulate struct{}
type Modulate1 struct {
	X0 Token
}

func (t Modulate) Apply(v Token) Token {
	return Modulate1{X0: v}
}

func (t Modulate1) Eval(c Context) Token {
	x0 := t.X0.Eval(c)
	r := Signal{S: mod(c, x0)}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Modulate) Eval(c Context) Token {
	return t
}

func (t Modulate) String() string {
	return "mod"
}

func (t Modulate1) String() string {
	return fmt.Sprintf("(mod %s)", t.X0)
}

type Demodulate struct{}
type Demodulate1 struct {
	X0 Token
}

func (t Demodulate) Apply(v Token) Token {
	return Demodulate1{X0: v}
}

func (t Demodulate1) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Signal).S
	r, s := demod(x0)
	if len(s) > 0 {
		log.Panicf("Invalid signal: %s", x0)
	}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Demodulate) Eval(c Context) Token {
	return t
}

func (t Demodulate) String() string {
	return "dem"
}

func (t Demodulate1) String() string {
	return fmt.Sprintf("(dem %s)", t.X0)
}

type Send struct{}

func (t Send) Eval(c Context) Token {
	return t
}

func (t Send) Apply(v Token) Token {
	log.Panicf("%s not implemented", t)
	return nil
}

func (t Send) String() string {
	return "send"
}

type Neg struct{}
type Neg1 struct {
	X0 Token
}

func (t Neg) Apply(v Token) Token {
	return Neg1{X0: v}
}

func (t Neg) Eval(c Context) Token {
	return t
}

func (t Neg1) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Value)
	r := Int{V: -x0.Value()}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Neg) String() string {
	return "neg"
}

func (t Neg1) String() string {
	return fmt.Sprintf("(neg %s)", t.X0)
}

type S struct{}
type S1 struct {
	X0 Token
}
type S2 struct {
	X0 Token
	X1 Token
}
type S3 struct {
	X0 Token
	X1 Token
	X2 Token
}

func (t S) Apply(v Token) Token {
	return S1{X0: v}
}

func (t S1) Apply(v Token) Token {
	return S2{X0: t.X0, X1: v}
}

func (t S2) Apply(v Token) Token {
	return S3{X0: t.X0, X1: t.X1, X2: v}
}

func (t S3) Eval(c Context) Token {
	y1 := t.X0.Eval(c).(Func)
	y2 := y1.Apply(t.X2).Eval(c).(Func)
	y3 := t.X1.(Func).Apply(t.X2)
	r := y2.Apply(y3).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t S) Eval(c Context) Token {
	return t
}

func (t S1) Eval(c Context) Token {
	return t
}

func (t S2) Eval(c Context) Token {
	return t
}

func (t S) String() string {
	return "s"
}

func (t S1) String() string {
	return fmt.Sprintf("(s %s)", t.X0)
}

func (t S2) String() string {
	return fmt.Sprintf("(s %s %s)", t.X0, t.X1)
}

func (t S3) String() string {
	return fmt.Sprintf("(s %s %s %s)", t.X0, t.X1, t.X2)
}

type C struct{}
type C1 struct {
	X0 Token
}
type C2 struct {
	X0 Token
	X1 Token
}
type C3 struct {
	X0 Token
	X1 Token
	X2 Token
}

func (t C) Apply(v Token) Token {
	return C1{X0: v}
}

func (t C1) Apply(v Token) Token {
	return C2{X0: t.X0, X1: v}
}

func (t C2) Apply(v Token) Token {
	return C3{X0: t.X0, X1: t.X1, X2: v}
}

func (t C3) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Func)
	y1 := x0.Apply(t.X2).Eval(c).(Func)
	r := y1.Apply(t.X1).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t C) Eval(c Context) Token {
	return t
}

func (t C1) Eval(c Context) Token {
	return t
}

func (t C2) Eval(c Context) Token {
	return t
}

func (t C) String() string {
	return "c"
}

func (t C1) String() string {
	return fmt.Sprintf("(c %s)", t.X0)
}

func (t C2) String() string {
	return fmt.Sprintf("(c %s %s)", t.X0, t.X1)
}

func (t C3) String() string {
	return fmt.Sprintf("(c %s %s %s)", t.X0, t.X1, t.X2)
}

type B struct{}
type B1 struct {
	X0 Token
}
type B2 struct {
	X0 Token
	X1 Token
}
type B3 struct {
	X0 Token
	X1 Token
	X2 Token
}

func (t B) Apply(v Token) Token {
	return B1{X0: v}
}

func (t B1) Apply(v Token) Token {
	return B2{X0: t.X0, X1: v}
}

func (t B2) Apply(v Token) Token {
	return B3{X0: t.X0, X1: t.X1, X2: v}
}

func (t B3) Eval(c Context) Token {
	r := Ap2{
		F: t.X0,
		A: Ap2{
			F: t.X1,
			A: t.X2,
		},
	}.Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t B) Eval(c Context) Token {
	return t
}

func (t B1) Eval(c Context) Token {
	return t
}

func (t B2) Eval(c Context) Token {
	return t
}

func (t B) String() string {
	return "b"
}

func (t B1) String() string {
	return fmt.Sprintf("(b %s)", t.X0)
}

func (t B2) String() string {
	return fmt.Sprintf("(b %s %s)", t.X0, t.X1)
}

func (t B3) String() string {
	return fmt.Sprintf("(b %s %s %s)", t.X0, t.X1, t.X2)
}

type Pwr2 struct{}
type Pwr21 struct {
	X0 Token
}

func (t Pwr2) Apply(v Token) Token {
	return Pwr21{X0: v}
}

func (t Pwr21) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Int).V
	r := Int{V: 1 << x0}
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Pwr2) Eval(c Context) Token {
	return t
}

func (t Pwr2) String() string {
	return "pwr2"
}

func (t Pwr21) String() string {
	return fmt.Sprintf("(pwr2 %s)", t.X0)
}

// type Pwr2 struct{}

// func (t Pwr2) Eval(c Context) Token {
// 	r := S2{
// 		X0: C2{
// 			X0: Eq1{
// 				X0: Int{V: 0},
// 			},
// 			X1: Int{V: 1},
// 		},
// 		X1: B2{
// 			X0: Mul1{
// 				X0: Int{V: 2},
// 			},
// 			X1: B2{
// 				X0: Pwr2{},
// 				X1: Dec{},
// 			},
// 		},
// 	}.Eval(c)
// 	// log.Printf("%s => %s", t, r)
// 	return r
// }

// func (t Pwr2) String() string {
// 	return "pwr2"
// }

type I struct{}
type I1 struct {
	X0 Token
}

func (t I) Apply(v Token) Token {
	return I1{X0: v}
}

func (t I1) Eval(c Context) Token {
	r := t.X0.Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t I) Eval(c Context) Token {
	return t
}

func (t I) String() string {
	return "i"
}

func (t I1) String() string {
	return fmt.Sprintf("(i %s)", t.X0)
}

type True struct{}
type True1 struct {
	X0 Token
}
type True2 struct {
	X0 Token
}

func (t True) Apply(v Token) Token {
	return True1{X0: v}
}

func (t True1) Apply(v Token) Token {
	return True2{X0: t.X0}
}

func (t True2) Eval(c Context) Token {
	r := t.X0.Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t True) Eval(c Context) Token {
	return t
}

func (t True1) Eval(c Context) Token {
	return t
}

func (t True) String() string {
	return "t"
}

func (t True1) String() string {
	return fmt.Sprintf("(t %s)", t.X0)
}

func (t True2) String() string {
	return fmt.Sprintf("(t %s ?)", t.X0)
}

type False struct{}
type False1 struct{}
type False2 struct {
	X1 Token
}

func (t False) Apply(v Token) Token {
	return False1{}
}

func (t False1) Apply(v Token) Token {
	return False2{X1: v}
}

func (t False2) Eval(c Context) Token {
	r := t.X1.Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t False) Eval(c Context) Token {
	return t
}

func (t False1) Eval(c Context) Token {
	return t
}

func (t False) String() string {
	return "f"
}

func (t False1) String() string {
	return "(f ?)"
}

func (t False2) String() string {
	return fmt.Sprintf("(f ? %s)", t.X1)
}

type Cons struct{}
type Cons1 struct {
	X0 Token
}
type Cons2 struct {
	X0 Token
	X1 Token
}
type Cons3 struct {
	X0 Token
	X1 Token
	X2 Token
}
type Vec struct {
	Cons
}

func (t Cons) Apply(v Token) Token {
	return Cons1{X0: v}
}

func (t Cons1) Apply(v Token) Token {
	return Cons2{X0: t.X0, X1: v}
}

func (t Cons2) Apply(v Token) Token {
	return Cons3{X0: t.X0, X1: t.X1, X2: v}
}

func (t Cons3) Eval(c Context) Token {
	y1 := t.X2.Eval(c).(Func).Apply(t.X0).Eval(c).(Func)
	r := y1.Apply(t.X1).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Cons) Eval(c Context) Token {
	return t
}

func (t Cons1) Eval(c Context) Token {
	return t
}

func (t Cons2) Eval(c Context) Token {
	return t
}

func (t Cons2) Car() Token {
	return t.X0
}

func (t Cons2) Cdr() Token {
	return t.X1
}

func (t Cons2) IsNil() bool {
	return false
}

func (t Cons) String() string {
	return "cons"
}

func (t Cons1) String() string {
	return fmt.Sprintf("(cons %s)", t.X0)
}

func (t Cons2) String() string {
	return fmt.Sprintf("(cons %s %s)", t.X0, t.X1)
}

func (t Cons3) String() string {
	return fmt.Sprintf("(cons %s %s %s)", t.X0, t.X1, t.X2)
}

type Car struct{}
type Car1 struct {
	X0 Token
}

func (t Car) Apply(v Token) Token {
	return Car1{X0: v}
}

func (t Car1) Eval(c Context) Token {
	y1 := t.X0.Eval(c).(Func)
	r := y1.Apply(True{}).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Car) Eval(c Context) Token {
	return t
}

func (t Car) String() string {
	return "car"
}

func (t Car1) String() string {
	return fmt.Sprintf("(car %s)", t.X0)
}

type Cdr struct{}
type Cdr1 struct {
	X0 Token
}

func (t Cdr) Apply(v Token) Token {
	return Cdr1{X0: v}
}

func (t Cdr1) Eval(c Context) Token {
	y1 := t.X0.Eval(c).(Func)
	r := y1.Apply(False{}).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Cdr) Eval(c Context) Token {
	return t
}

func (t Cdr) String() string {
	return "cdr"
}

func (t Cdr1) String() string {
	return fmt.Sprintf("(cdr %s)", t.X0)
}

type Nil struct{}
type Nil1 struct{}

func (t Nil) Apply(v Token) Token {
	return Nil1{}
}

func (t Nil1) Eval(c Context) Token {
	return True{}
}

func (t Nil) Eval(c Context) Token {
	return t
}

func (t Nil) Car() Token {
	return Nil{}
}

func (t Nil) Cdr() Token {
	return Nil{}
}

func (t Nil) IsNil() bool {
	return true
}

func (t Nil) String() string {
	return "nil"
}

func (t Nil1) String() string {
	return "(nil ?)"
}

type isNil struct{}
type isNil1 struct{}
type isNil2 struct{}

func (t isNil) Apply(v Token) Token {
	return isNil1{}
}

func (t isNil1) Apply(v Token) Token {
	return isNil2{}
}

func (t isNil2) Eval(c Context) Token {
	return False{}
}

func (t isNil) Eval(c Context) Token {
	return t
}

func (t isNil1) Eval(c Context) Token {
	return t
}

func (t isNil) String() string {
	return "*isnil*"
}

func (t isNil1) String() string {
	return "(*isnil* ?)"
}

func (t isNil2) String() string {
	return "(*isnil* ? ?)"
}

type IsNil struct{}
type IsNil1 struct {
	X0 Token
}

func (t IsNil) Apply(v Token) Token {
	return IsNil1{X0: v}
}

func (t IsNil1) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Func)
	r := x0.Apply(isNil{}).Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t IsNil) Eval(c Context) Token {
	return t
}

func (t IsNil) String() string {
	return "isnil"
}

func (t IsNil1) String() string {
	return fmt.Sprintf("(isnil %s)", t.X0)
}

type Point struct {
	X int64
	Y int64
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

type Picture []Point

func (ps Picture) Eval(c Context) Token {
	return ps
}

func (ps Picture) String() string {
	var pp []string
	for _, p := range ps {
		pp = append(pp, p.String())
	}
	return "(" + strings.Join(pp, " ") + ")"
}

type ICons interface {
	Token
	Car() Token
	Cdr() Token
	IsNil() bool
}

func ListPoints(c Context, v Token) Picture {
	var r Picture
	for i := v.(ICons); !i.IsNil(); i = i.Cdr().Eval(c).(ICons) {
		p := i.Car().Eval(c).(ICons)
		x := p.Car().Eval(c).(Int).V
		y := p.Cdr().Eval(c).(Int).V
		r = append(r, Point{X: x, Y: y})
	}
	return r
}

func DrawPoints(c Context, v Token) Picture {
	ps := ListPoints(c, v.Eval(c))
	log.Printf("Draw %s", ps)
	return ps
}

type Draw struct{}
type Draw1 struct {
	X0 Token
}

func (t Draw) Apply(v Token) Token {
	return Draw1{X0: v}
}

func (t Draw1) Eval(c Context) Token {
	return DrawPoints(c, t.X0)
}

func (t Draw) Eval(c Context) Token {
	return t
}

func (t Draw) String() string {
	return "draw"
}

func (t Draw1) String() string {
	return fmt.Sprintf("(draw %s)", t.X0)
}

type Checkerboard struct{}

func (t Checkerboard) Apply(v Token) Token {
	log.Panicf("%s not implemented", t)
	return nil
}

func (t Checkerboard) Eval(c Context) Token {
	return t
}

func (t Checkerboard) String() string {
	return "checkerboard"
}

type Multipledraw struct{}
type Multipledraw1 struct {
	X0 Token
}

func (t Multipledraw) Apply(v Token) Token {
	return Multipledraw1{X0: v}
}

func (t Multipledraw1) Eval(c Context) Token {
	var r Picture
	v := t.X0.Eval(c)
	for i := v.(ICons); !i.IsNil(); i = i.Cdr().Eval(c).(ICons) {
		r = append(r, DrawPoints(c, i.Car())...)
	}
	return r
}

func (t Multipledraw) Eval(c Context) Token {
	return t
}

func (t Multipledraw) String() string {
	return "multipledraw"
}

func (t Multipledraw1) String() string {
	return fmt.Sprintf("(multipledraw %s)", t.X0)
}

type If0 struct{}
type If01 struct {
	X0 Token
}
type If02 struct {
	X0 Token
	X1 Token
}
type If03 struct {
	X0 Token
	X1 Token
	X2 Token
}

func (t If0) Apply(v Token) Token {
	return If01{X0: v}
}

func (t If01) Apply(v Token) Token {
	return If02{X0: t.X0, X1: v}
}

func (t If02) Apply(v Token) Token {
	return If03{X0: t.X0, X1: t.X1, X2: v}
}

func (t If03) Eval(c Context) Token {
	x0 := t.X0.Eval(c).(Int).V
	r := t.X2
	if x0 == 0 {
		r = t.X1
	}
	// r = r.Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t If0) Eval(c Context) Token {
	return t
}

func (t If01) Eval(c Context) Token {
	return t
}

func (t If02) Eval(c Context) Token {
	return t
}

func (t If0) String() string {
	return "if0"
}

func (t If01) String() string {
	return fmt.Sprintf("(if0 %s)", t.X0)
}

func (t If02) String() string {
	return fmt.Sprintf("(if0 %s %s)", t.X0, t.X1)
}

func (t If03) String() string {
	return fmt.Sprintf("(if0 %s %s %s)", t.X0, t.X1, t.X2)
}
