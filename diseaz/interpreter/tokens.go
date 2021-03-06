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
	Galaxy() string
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

func (t Ap) Galaxy() string {
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
	f := c.Eval(t.F).(Func)
	r := f.Apply(t.A)
	// log.Printf("%s => %s", t, r)
	return r, true
}

func (t Ap2) Galaxy() string {
	return fmt.Sprintf("ap %s %s", t.F.Galaxy(), t.A.Galaxy())
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

func (t VarN) Galaxy() string {
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

func (t Int) Galaxy() string {
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
	v := c.Eval(t.X0).(Int)
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

func (t Inc) Galaxy() string {
	return "inc"
}

func (t Inc1) String() string {
	return fmt.Sprintf("(inc %s)", t.X0)
}

func (t Inc1) Galaxy() string {
	return fmt.Sprintf("ap inc %s", t.X0.Galaxy())
}

type Dec struct{}
type Dec1 struct {
	X0 Token
}

func (t Dec) Apply(v Token) Token {
	return Dec1{X0: v}
}

func (t Dec1) Eval(c Context) (Token, bool) {
	v := c.Eval(t.X0).(Int)
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

func (t Dec) Galaxy() string {
	return "dec"
}

func (t Dec1) Galaxy() string {
	return fmt.Sprintf("ap dec %s", t.X0.Galaxy())
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
	x0 := c.Eval(t.X0).(Int)
	x1 := c.Eval(t.X1).(Int)
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

func (t Add) Galaxy() string {
	return "add"
}

func (t Add1) Galaxy() string {
	return fmt.Sprintf("ap add %s", t.X0.Galaxy())
}

func (t Add2) Galaxy() string {
	return fmt.Sprintf("ap ap add %s %s", t.X0.Galaxy(), t.X1.Galaxy())
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
	x0 := c.Eval(t.X0).(Int)
	x1 := c.Eval(t.X1).(Int)
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

func (t Mul) Galaxy() string {
	return "mul"
}

func (t Mul1) Galaxy() string {
	return fmt.Sprintf("ap mul %s", t.X0.Galaxy())
}

func (t Mul2) Galaxy() string {
	return fmt.Sprintf("ap ap mul %s %s", t.X0.Galaxy(), t.X1.Galaxy())
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
	x0 := c.Eval(t.X0).(Int)
	x1 := c.Eval(t.X1).(Int)
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

func (t Div) Galaxy() string {
	return "div"
}

func (t Div1) Galaxy() string {
	return fmt.Sprintf("ap div %s", t.X0.Galaxy())
}

func (t Div2) Galaxy() string {
	return fmt.Sprintf("ap ap div %s %s", t.X0.Galaxy(), t.X1.Galaxy())
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
	x0 := c.Eval(t.X0).(Int)
	x1 := c.Eval(t.X1).(Int)
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

func (t Eq) Galaxy() string {
	return "eq"
}

func (t Eq1) Galaxy() string {
	return fmt.Sprintf("ap eq %s", t.X0.Galaxy())
}

func (t Eq2) Galaxy() string {
	return fmt.Sprintf("ap ap eq %s %s", t.X0.Galaxy(), t.X1.Galaxy())
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
	x0 := c.Eval(t.X0).(Int)
	x1 := c.Eval(t.X1).(Int)
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

func (t Lt) Galaxy() string {
	return "lt"
}

func (t Lt1) Galaxy() string {
	return fmt.Sprintf("ap lt %s", t.X0.Galaxy())
}

func (t Lt2) Galaxy() string {
	return fmt.Sprintf("ap ap lt %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func ModulateToken(v Token) string {
	switch tt := v.(type) { // TailEval?
	case Int:
		return modInt(tt.V)
	case ICons:
		return modCons(tt)
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

func modCons(v ICons) string {
	if v.IsNil() {
		return "00"
	}
	return "11" + ModulateToken(v.Car()) + ModulateToken(v.Cdr())
}

func DemodulateToken(v string) Token {
	r, s := demodToken(v)
	if len(s) > 0 {
		log.Panicf("Extra tail on demod %s\n=> %s\n++ %s", v, r, s)
	}
	return r
}

func demodToken(v string) (Token, string) {
	if len(v) == 0 {
		return nil, v
	}
	prefix, w := v[0:2], v[2:]
	if prefix == "00" {
		return Nil{}, w
	}
	if prefix == "11" {
		car, w := demodToken(w)
		cdr, w := demodToken(w)
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

func (t Signal) Galaxy() string {
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
	x0 := c.Eval(t.X0)
	r := Signal{S: ModulateToken(x0)}
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

func (t Modulate) Galaxy() string {
	return "mod"
}

func (t Modulate1) Galaxy() string {
	return fmt.Sprintf("ap mod %s", t.X0.Galaxy())
}

type Demodulate struct{}
type Demodulate1 struct {
	X0 Token
}

func (t Demodulate) Apply(v Token) Token {
	return Demodulate1{X0: v}
}

func (t Demodulate1) Eval(c Context) (Token, bool) {
	x0 := c.Eval(t.X0).(Signal).S
	r := DemodulateToken(x0)
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

func (t Demodulate) Galaxy() string {
	return "dem"
}

func (t Demodulate1) Galaxy() string {
	return fmt.Sprintf("ap dem %s", t.X0.Galaxy())
}

type Send struct{}
type Send1 struct {
	X0 Token
}

func (t Send) Apply(v Token) Token {
	return Send1{X0: v}
}

func (t Send1) Eval(c Context) (Token, bool) {
	x0 := c.Eval(t.X0)
	log.Printf("Sending: %s", x0)
	r := DemodulateToken(c.Send(ModulateToken(x0)))
	log.Printf("Receivd: %s", r.Galaxy())
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

func (t Send) Galaxy() string {
	return "send"
}

func (t Send1) Galaxy() string {
	return fmt.Sprintf(" ap send %s", t.X0.Galaxy())
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
	x0 := c.Eval(t.X0).(Int)
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

func (t Neg) Galaxy() string {
	return "neg"
}

func (t Neg1) Galaxy() string {
	return fmt.Sprintf("ap neg %s", t.X0.Galaxy())
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
	f0 := c.Eval(t.X0).(Func)
	x2 := c.Eval(t.X2)
	f1 := c.Eval(f0.Apply(x2)).(Func)
	r := f1.Apply(
		Ap2{
			F: t.X1,
			A: x2,
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

func (t S) Galaxy() string {
	return "s"
}

func (t S1) Galaxy() string {
	return fmt.Sprintf("ap s %s", t.X0.Galaxy())
}

func (t S2) Galaxy() string {
	return fmt.Sprintf("ap ap s %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func (t S3) Galaxy() string {
	return fmt.Sprintf("ap ap ap s %s %s %s", t.X0.Galaxy(), t.X1.Galaxy(), t.X2.Galaxy())
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
	f0 := c.Eval(t.X0).(Func)
	f1 := c.Eval(f0.Apply(t.X2)).(Func)
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

func (t C) Galaxy() string {
	return "c"
}

func (t C1) Galaxy() string {
	return fmt.Sprintf("ap c %s", t.X0.Galaxy())
}

func (t C2) Galaxy() string {
	return fmt.Sprintf("ap ap c %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func (t C3) Galaxy() string {
	return fmt.Sprintf("ap ap ap c %s %s %s", t.X0.Galaxy(), t.X1.Galaxy(), t.X2.Galaxy())
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
	r := c.Eval(t.X0).(Func).Apply(
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

func (t B) Galaxy() string {
	return "b"
}

func (t B1) Galaxy() string {
	return fmt.Sprintf("ap b %s", t.X0.Galaxy())
}

func (t B2) Galaxy() string {
	return fmt.Sprintf("ap ap b %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func (t B3) Galaxy() string {
	return fmt.Sprintf("ap ap ap b %s %s %s", t.X0.Galaxy(), t.X1.Galaxy(), t.X2.Galaxy())
}

type Pwr2 struct{}
type Pwr21 struct {
	X0 Token
}

func (t Pwr2) Apply(v Token) Token {
	return Pwr21{X0: v}
}

func (t Pwr21) Eval(c Context) (Token, bool) {
	x0 := c.Eval(t.X0).(Int).V
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

func (t Pwr2) Galaxy() string {
	return "pwr2"
}

func (t Pwr21) Galaxy() string {
	return fmt.Sprintf("ap pwr2 %s", t.X0.Galaxy())
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

func (t I) Galaxy() string {
	return "i"
}

func (t I1) Galaxy() string {
	return fmt.Sprintf("ap i %s", t.X0.Galaxy())
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

func (t True) Galaxy() string {
	return "t"
}

func (t True1) Galaxy() string {
	return fmt.Sprintf("ap t %s", t.X0.Galaxy())
}

func (t True2) Galaxy() string {
	return fmt.Sprintf("ap ap t %s 42", t.X0.Galaxy())
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

func (t False) Galaxy() string {
	return "f"
}

func (t False1) Galaxy() string {
	return "ap f 42"
}

func (t False2) Galaxy() string {
	return fmt.Sprintf("ap ap f 42 %s", t.X1.Galaxy())
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
	y1 := c.Eval(t.X2).(Func)
	y2 := c.Eval(y1.Apply(t.X0)).(Func)
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
	t.X0 = c.Eval(t.X0)
	t.X1 = c.Eval(t.X1)
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

func (t Cons) Galaxy() string {
	return "cons"
}

func (t Cons1) Galaxy() string {
	return fmt.Sprintf("ap cons %s", t.X0.Galaxy())
}

func (t Cons2) Galaxy() string {
	return fmt.Sprintf("ap ap cons %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func (t Cons3) Galaxy() string {
	return fmt.Sprintf("ap ap ap cons %s %s %s", t.X0.Galaxy(), t.X1.Galaxy(), t.X2.Galaxy())
}

type Vec struct {
	Cons
}

func (t Vec) String() string {
	return "vec"
}

func (t Vec) Galaxy() string {
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
	y1 := c.Eval(t.X0).(Func)
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

func (t Car) Galaxy() string {
	return "car"
}

func (t Car1) Galaxy() string {
	return fmt.Sprintf("ap car %s", t.X0.Galaxy())
}

type Cdr struct{}
type Cdr1 struct {
	X0 Token
}

func (t Cdr) Apply(v Token) Token {
	return Cdr1{X0: v}
}

func (t Cdr1) Eval(c Context) (Token, bool) {
	y1 := c.Eval(t.X0).(Func)
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

func (t Cdr) Galaxy() string {
	return "cdr"
}

func (t Cdr1) Galaxy() string {
	return fmt.Sprintf("ap cdr %s", t.X0.Galaxy())
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

func (t Nil) Galaxy() string {
	return "nil"
}

func (t Nil1) Galaxy() string {
	return "ap nil 42"
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

func (t isNil) Galaxy() string {
	return "*isnil*"
}

func (t isNil1) Galaxy() string {
	return "ap *isnil* 42"
}

func (t isNil2) Galaxy() string {
	return "ap ap *isnil* 42 42"
}

type IsNil struct{}
type IsNil1 struct {
	X0 Token
}

func (t IsNil) Apply(v Token) Token {
	return IsNil1{X0: v}
}

func (t IsNil1) Eval(c Context) (Token, bool) {
	x0 := c.Eval(t.X0).(Func)
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

func (t IsNil) Galaxy() string {
	return "isnil"
}

func (t IsNil1) Galaxy() string {
	return fmt.Sprintf("ap isnil %s", t.X0.Galaxy())
}

type ICons interface {
	Token
	Car() Token
	Cdr() Token
	IsNil() bool
}

func DrawPoints(c Context, v Token) *Picture {
	var pts []Point
	r := NewPicture()
	for i := c.Eval(v).(ICons); !i.IsNil(); i = i.Cdr().(ICons) {
		p := i.Car().(ICons)
		x := int(p.Car().(Int).V)
		y := int(p.Cdr().(Int).V)
		pts = append(pts, Pt(x, y))
	}
	c.Picture().DrawPts(pts...)
	r.DrawPts(pts...)
	// log.Printf("Draw %s", r)
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

func (t Draw) Galaxy() string {
	return "draw"
}

func (t Draw1) Galaxy() string {
	return fmt.Sprintf("ap draw %s", t.X0.Galaxy())
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

func (t Checkerboard) Galaxy() string {
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
	r := NewPicture()
	v := c.Eval(t.X0).(ICons)
	for i := v; !i.IsNil(); i = c.Eval(i.Cdr()).(ICons) {
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

func (t Multipledraw) Galaxy() string {
	return "multipledraw"
}

func (t Multipledraw1) Galaxy() string {
	return fmt.Sprintf("ap multipledraw %s", t.X0.Galaxy())
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
	x0 := c.Eval(t.X0).(Int).V
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

func (t If0) Galaxy() string {
	return "if0"
}

func (t If01) Galaxy() string {
	return fmt.Sprintf("ap if0 %s", t.X0.Galaxy())
}

func (t If02) Galaxy() string {
	return fmt.Sprintf("ap ap if0 %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func (t If03) Galaxy() string {
	return fmt.Sprintf("ap ap ap if0 %s %s %s", t.X0.Galaxy(), t.X1.Galaxy(), t.X2.Galaxy())
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

	x1 := c.Eval(t.X1).(ICons)
	flag := c.Eval(x1.Car()).(Int)
	x11 := c.Eval(x1.Cdr()).(ICons)
	newState := x11.Car()
	x12 := c.Eval(x11.Cdr()).(ICons)
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

func (t interactHelper) Galaxy() string {
	return "*interact-helper*"
}

func (t interactHelper1) Galaxy() string {
	return fmt.Sprintf("ap *interact-helper* %s", t.X0.Galaxy())
}

func (t interactHelper2) Galaxy() string {
	return fmt.Sprintf("ap ap *interact-helper* %s %s", t.X0.Galaxy(), t.X1.Galaxy())
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

func (t Interact) Galaxy() string {
	return "interact"
}

func (t Interact1) Galaxy() string {
	return fmt.Sprintf("ap interact %s", t.X0.Galaxy())
}

func (t Interact2) Galaxy() string {
	return fmt.Sprintf("ap ap interact %s %s", t.X0.Galaxy(), t.X1.Galaxy())
}

func (t Interact3) Galaxy() string {
	return fmt.Sprintf("ap ap ap interact %s %s %s", t.X0.Galaxy(), t.X1.Galaxy(), t.X2.Galaxy())
}
