package interpreter

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Token interface {
	Eval(c Context) (Token, bool)
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

func (t Ap) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Ap2) Eval(c Context) (Token, bool) {
	f := TailEval(c, t.F).(Func)
	r := f.Apply(t.A)
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t Ap2) String() string {
	return fmt.Sprintf("(ap %s %s)", t.F, t.A)
}

type VarN struct {
	N int
}

func (t VarN) Eval(c Context) (Token, bool) {
	return c.GetVar(t.N), true
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

func (t Int) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Inc1) Eval(c Context) (Token, bool) {
	v := TailEval(c, t.X0).(Int)
	r := Int{V: v.V + 1}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Inc) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Dec1) Eval(c Context) (Token, bool) {
	v := TailEval(c, t.X0).(Int)
	r := Int{V: v.V - 1}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Dec) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Add2) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int)
	x1 := TailEval(c, t.X1).(Int)
	r := Int{V: x0.V + x1.V}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Add) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Add1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Add) String() string {
	return "add"
}

func (t Add1) String() string {
	return fmt.Sprintf("(add1 %s)", t.X0)
}

func (t Add2) String() string {
	return fmt.Sprintf("(add2 %s %s)", t.X0, t.X1)
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

func (t Mul2) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int)
	x1 := TailEval(c, t.X1).(Int)
	r := Int{V: x0.V * x1.V}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Mul) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Mul1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Mul) String() string {
	return "mul"
}

func (t Mul1) String() string {
	return fmt.Sprintf("(mul1 %s)", t.X0)
}

func (t Mul2) String() string {
	return fmt.Sprintf("(mul2 %s %s)", t.X0, t.X1)
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

func (t Div2) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int)
	x1 := TailEval(c, t.X1).(Int)
	r := Int{V: x0.V / x1.V}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Div) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Div1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Div) String() string {
	return "div"
}

func (t Div1) String() string {
	return fmt.Sprintf("(div1 %s)", t.X0)
}

func (t Div2) String() string {
	return fmt.Sprintf("(div2 %s %s)", t.X0, t.X1)
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

func (t Eq2) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int)
	x1 := TailEval(c, t.X1).(Int)
	var r Token = False{}
	if x0.V == x1.V {
		r = True{}
	}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Eq) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Eq1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Eq) String() string {
	return "eq"
}

func (t Eq1) String() string {
	return fmt.Sprintf("(eq1 %s)", t.X0)
}

func (t Eq2) String() string {
	return fmt.Sprintf("(eq2 %s %s)", t.X0, t.X1)
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

func (t Lt2) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int)
	x1 := TailEval(c, t.X1).(Int)
	var r Token = False{}
	if x0.V < x1.V {
		r = True{}
	}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Lt) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Lt1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Lt) String() string {
	return "lt"
}

func (t Lt1) String() string {
	return fmt.Sprintf("(lt1 %s)", t.X0)
}

func (t Lt2) String() string {
	return fmt.Sprintf("(lt2 %s %s)", t.X0, t.X1)
}

func mod(c Context, v Token) string {
	switch tt := v.(type) { // TailEval?
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

func (t Signal) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Modulate1) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0)
	r := Signal{S: mod(c, x0)}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Modulate) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Demodulate1) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Signal).S
	r, s := demod(x0)
	if len(s) > 0 {
		log.Panicf("Invalid signal: %s", x0)
	}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Demodulate) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Demodulate) String() string {
	return "dem"
}

func (t Demodulate1) String() string {
	return fmt.Sprintf("(dem %s)", t.X0)
}

type Send struct{}
type Send1 struct {
	X0 Token
}

func (t Send) Apply(v Token) Token {
	return Send1{X0: v}
}

func (t Send1) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0)
	rm := c.Send(mod(c, x0))
	r, s := demod(rm)
	if len(s) > 0 {
		log.Panicf("Extra tail on demod %#v = %#v", rm, s)
	}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Send) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Send) String() string {
	return "send"
}

func (t Send1) String() string {
	return fmt.Sprintf("(send %s)", t.X0)
}

type Neg struct{}
type Neg1 struct {
	X0 Token
}

func (t Neg) Apply(v Token) Token {
	return Neg1{X0: v}
}

func (t Neg) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Neg1) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int)
	r := Int{V: -x0.V}
	// log.Printf("%s => %s", t, r)
	return r, false
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

func (t S3) Eval(c Context) (Token, bool) {
	f0 := TailEval(c, t.X0).(Func)
	f1 := TailEval(c, f0.Apply(t.X2)).(Func)
	r := f1.Apply(
		Ap2{
			F: t.X1,
			A: t.X2,
		},
	)
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t S) Eval(c Context) (Token, bool) {
	return t, false
}

func (t S1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t S2) Eval(c Context) (Token, bool) {
	return t, false
}

func (t S) String() string {
	return "s"
}

func (t S1) String() string {
	return fmt.Sprintf("(s1 %s)", t.X0)
}

func (t S2) String() string {
	return fmt.Sprintf("(s2 %s %s)", t.X0, t.X1)
}

func (t S3) String() string {
	return fmt.Sprintf("(s3 %s %s %s)", t.X0, t.X1, t.X2)
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

func (t C3) Eval(c Context) (Token, bool) {
	f0 := TailEval(c, t.X0).(Func)
	f1 := TailEval(c, f0.Apply(t.X2)).(Func)
	r := f1.Apply(t.X1)
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t C) Eval(c Context) (Token, bool) {
	return t, false
}

func (t C1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t C2) Eval(c Context) (Token, bool) {
	return t, false
}

func (t C) String() string {
	return "c"
}

func (t C1) String() string {
	return fmt.Sprintf("(c1 %s)", t.X0)
}

func (t C2) String() string {
	return fmt.Sprintf("(c2 %s %s)", t.X0, t.X1)
}

func (t C3) String() string {
	return fmt.Sprintf("(c3 %s %s %s)", t.X0, t.X1, t.X2)
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

func (t B3) Eval(c Context) (Token, bool) {
	r := TailEval(c, t.X0).(Func).Apply(
		Ap2{
			F: t.X1,
			A: t.X2,
		},
	)
	return r, true
}

func (t B) Eval(c Context) (Token, bool) {
	return t, false
}

func (t B1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t B2) Eval(c Context) (Token, bool) {
	return t, false
}

func (t B) String() string {
	return "b"
}

func (t B1) String() string {
	return fmt.Sprintf("(b1 %s)", t.X0)
}

func (t B2) String() string {
	return fmt.Sprintf("(b2 %s %s)", t.X0, t.X1)
}

func (t B3) String() string {
	return fmt.Sprintf("(b3 %s %s %s)", t.X0, t.X1, t.X2)
}

type Pwr2 struct{}
type Pwr21 struct {
	X0 Token
}

func (t Pwr2) Apply(v Token) Token {
	return Pwr21{X0: v}
}

func (t Pwr21) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int).V
	r := Int{V: 1 << x0}
	// log.Printf("%s => %s", t, r)
	return r, false
}

func (t Pwr2) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Pwr2) String() string {
	return "pwr2"
}

func (t Pwr21) String() string {
	return fmt.Sprintf("(pwr2 %s)", t.X0)
}

type I struct{}
type I1 struct {
	X0 Token
}

func (t I) Apply(v Token) Token {
	return I1{X0: v}
}

func (t I1) Eval(c Context) (Token, bool) {
	r := t.X0
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t I) Eval(c Context) (Token, bool) {
	return t, false
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

func (t True2) Eval(c Context) (Token, bool) {
	r := t.X0
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t True) Eval(c Context) (Token, bool) {
	return t, false
}

func (t True1) Eval(c Context) (Token, bool) {
	return t, false
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

func (t False2) Eval(c Context) (Token, bool) {
	r := t.X1
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t False) Eval(c Context) (Token, bool) {
	return t, false
}

func (t False1) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Cons) Apply(v Token) Token {
	return Cons1{X0: v}
}

func (t Cons1) Apply(v Token) Token {
	return Cons2{X0: t.X0, X1: v}
}

func (t Cons2) Apply(v Token) Token {
	return Cons3{X0: t.X0, X1: t.X1, X2: v}
}

func (t Cons3) Eval(c Context) (Token, bool) {
	y1 := TailEval(c, t.X2).(Func)
	y2 := TailEval(c, y1.Apply(t.X0)).(Func)
	r := y2.Apply(t.X1)
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t Cons) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Cons1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Cons2) Eval(c Context) (Token, bool) {
	t.X0 = TailEval(c, t.X0)
	t.X1 = TailEval(c, t.X1)
	return t, false
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
	return fmt.Sprintf("(cons1 %s)", t.X0)
}

func (t Cons2) String() string {
	return fmt.Sprintf("(cons2 %s %s)", t.X0, t.X1)
}

func (t Cons3) String() string {
	return fmt.Sprintf("(cons3 %s %s %s)", t.X0, t.X1, t.X2)
}

type Vec struct {
	Cons
}

func (t Vec) String() string {
	return "vec"
}

type Car struct{}
type Car1 struct {
	X0 Token
}

func (t Car) Apply(v Token) Token {
	return Car1{X0: v}
}

func (t Car1) Eval(c Context) (Token, bool) {
	y1 := TailEval(c, t.X0).(Func)
	r := y1.Apply(True{})
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t Car) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Cdr1) Eval(c Context) (Token, bool) {
	y1 := TailEval(c, t.X0).(Func)
	r := y1.Apply(False{})
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t Cdr) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Nil1) Eval(c Context) (Token, bool) {
	return True{}, false
}

func (t Nil) Eval(c Context) (Token, bool) {
	return t, false
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

func (t isNil2) Eval(c Context) (Token, bool) {
	return False{}, false
}

func (t isNil) Eval(c Context) (Token, bool) {
	return t, false
}

func (t isNil1) Eval(c Context) (Token, bool) {
	return t, false
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

func (t IsNil1) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Func)
	r := x0.Apply(isNil{})
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t IsNil) Eval(c Context) (Token, bool) {
	return t, false
}

func (t IsNil) String() string {
	return "isnil"
}

func (t IsNil1) String() string {
	return fmt.Sprintf("(isnil %s)", t.X0)
}

type ICons interface {
	Token
	Car() Token
	Cdr() Token
	IsNil() bool
}

// func ListPoints(c Context, v Token) Picture {
// 	var r Picture
// 	for i := v.(ICons); !i.IsNil(); i = TailEval(c, i.Cdr()).(ICons) {
// 		p := TailEval(c, i.Car()).(ICons)
// 		x := TailEval(c, p.Car()).(Int).V
// 		y := TailEval(c, p.Cdr()).(Int).V
// 		r = append(r, Point{X: x, Y: y})
// 	}
// 	return r
// }

func DrawPoints(c Context, v Token) Picture {
	pic := c.Picture()
	r := Picture{}
	for i := TailEval(c, v).(ICons); !i.IsNil(); i = TailEval(c, i.Cdr()).(ICons) {
		p := TailEval(c, i.Car()).(ICons)
		x := int(TailEval(c, p.Car()).(Int).V)
		y := int(TailEval(c, p.Cdr()).(Int).V)
		pic.Draw(x, y)
		r.Draw(x, y)
	}
	log.Printf("Draw %s", r)
	return r
}

type Draw struct{}
type Draw1 struct {
	X0 Token
}

func (t Draw) Apply(v Token) Token {
	return Draw1{X0: v}
}

func (t Draw1) Eval(c Context) (Token, bool) {
	return DrawPoints(c, t.X0), false
}

func (t Draw) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Checkerboard) Eval(c Context) (Token, bool) {
	return t, false
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

func (t Multipledraw1) Eval(c Context) (Token, bool) {
	r := Picture{}
	v := TailEval(c, t.X0).(ICons)
	for i := v; !i.IsNil(); i = TailEval(c, i.Cdr()).(ICons) {
		r.DrawPicture(DrawPoints(c, i.Car()))
	}
	return r, false
}

func (t Multipledraw) Eval(c Context) (Token, bool) {
	return t, false
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

func (t If03) Eval(c Context) (Token, bool) {
	x0 := TailEval(c, t.X0).(Int).V
	r := t.X2
	if x0 == 0 {
		r = t.X1
	}
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t If0) Eval(c Context) (Token, bool) {
	return t, false
}

func (t If01) Eval(c Context) (Token, bool) {
	return t, false
}

func (t If02) Eval(c Context) (Token, bool) {
	return t, false
}

func (t If0) String() string {
	return "if0"
}

func (t If01) String() string {
	return fmt.Sprintf("(if0_1 %s)", t.X0)
}

func (t If02) String() string {
	return fmt.Sprintf("(if0_2 %s %s)", t.X0, t.X1)
}

func (t If03) String() string {
	return fmt.Sprintf("(if0_3 %s %s %s)", t.X0, t.X1, t.X2)
}

type interactHelper struct{}
type interactHelper1 struct {
	X0 Token
}
type interactHelper2 struct {
	X0 Token
	X1 Token
}

func (t interactHelper) Apply(v Token) Token {
	return interactHelper1{X0: v}
}

func (t interactHelper1) Apply(v Token) Token {
	return interactHelper2{X0: t.X0, X1: v}
}

func (t interactHelper2) Eval(c Context) (Token, bool) {
	// (f38 x0 x1) =
	//   (if0 (car x1)
	// 	  -- then
	//    (cons (modem (car (cdr x1))) (cons (multipledraw (car (cdr (cdr x1)))) nil))
	// 	  -- else
	// 	  (interact x0 (modem (car (cdr x1))) (send (car (cdr (cdr x1))))))

	// f38(protocol, (flag, newState, data)) = if flag == 0
	//                 then (modem(newState), multipledraw(data))
	//                 else interact(protocol, modem(newState), send(data))

	x1 := TailEval(c, t.X1).(ICons)
	flag := TailEval(c, x1.Car()).(Int)
	x11 := TailEval(c, x1.Cdr()).(ICons)
	newState := x11.Car()
	x12 := TailEval(c, x11.Cdr()).(ICons)
	data := x12.Car()
	if flag.V == 0 {
		r := Cons2{
			X0: newState,
			X1: Cons2{
				X0: Multipledraw1{X0: data},
				X1: Nil{},
			},
		}
		// log.Printf("%s => %s", t, r)
		return r, true
	} else {
		r := Interact3{
			X0: t.X0,
			X1: newState,
			X2: Send1{
				X0: data,
			},
		}
		// log.Printf("%s => %s", t, r)
		return r, true
	}
}

func (t interactHelper) Eval(c Context) (Token, bool) {
	return t, false
}

func (t interactHelper1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t interactHelper) String() string {
	return "interact-helper"
}

func (t interactHelper1) String() string {
	return fmt.Sprintf("(interact-helper1 %s)", t.X0)
}

func (t interactHelper2) String() string {
	return fmt.Sprintf("(interact-helper2 %s %s)", t.X0, t.X1)
}

type Interact struct{}
type Interact1 struct {
	X0 Token
}
type Interact2 struct {
	X0 Token
	X1 Token
}
type Interact3 struct {
	X0 Token
	X1 Token
	X2 Token
}

func (t Interact) Apply(v Token) Token {
	return Interact1{X0: v}
}

func (t Interact1) Apply(v Token) Token {
	return Interact2{X0: t.X0, X1: v}
}

func (t Interact2) Apply(v Token) Token {
	return Interact3{X0: t.X0, X1: t.X1, X2: v}
}

func (t Interact3) Eval(c Context) (Token, bool) {
	// (interact x0 x1 x2) = (f38 x0 (x0 x1 x2))

	r := interactHelper2{
		X0: t.X0,
		X1: Ap2{
			F: Ap2{
				F: t.X0,
				A: t.X1,
			},
			A: t.X2,
		},
	}
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t Interact) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Interact1) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Interact2) Eval(c Context) (Token, bool) {
	return t, false
}

func (t Interact) String() string {
	return "interact"
}

func (t Interact1) String() string {
	return fmt.Sprintf("(interact1 %s)", t.X0)
}

func (t Interact2) String() string {
	return fmt.Sprintf("(interact2 %s %s)", t.X0, t.X1)
}

func (t Interact3) String() string {
	return fmt.Sprintf("(interact3 %s %s %s)", t.X0, t.X1, t.X2)
}
