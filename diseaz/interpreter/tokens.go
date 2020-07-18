package interpreter

import (
	"fmt"
	"strconv"
)

type Token interface {
	Eval(c Context) Token
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

// type Modulate struct{}

// func (t Modulate) Apply(v Token) (Program, error) {
// 	return nil, fmt.Errorf("%s not implemented", t.String())
// }

// func (t Modulate) String() string {
// 	return "mod"
// }

// type Demodulate struct{}

// func (t Demodulate) Apply(v Token) (Program, error) {
// 	return nil, fmt.Errorf("%s not implemented", t.String())
// }

// func (t Demodulate) String() string {
// 	return "dem"
// }

// type Send struct{}

// func (t Send) Apply(v Token) (Program, error) {
// 	return nil, fmt.Errorf("%s not implemented", t.String())
// }

// func (t Send) String() string {
// 	return "send"
// }

// type Neg struct{}

// func (t Neg) Apply(v Token) (Program, error) {
// 	w, ok := v.(Value)
// 	if !ok {
// 		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
// 	}
// 	return Program{Int{V: -w.Value()}}, nil
// }

// func (t Neg) String() string {
// 	return "neg"
// }

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
	r := Ap2{
		F: Ap2{
			F: t.X0,
			A: t.X2,
		},
		A: Ap2{
			F: t.X1,
			A: t.X2,
		},
	}.Eval(c)
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
	r := Ap2{
		F: Ap2{
			F: t.X0,
			A: t.X2,
		},
		A: t.X1,
	}.Eval(c)
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

func (t Pwr2) Eval(c Context) Token {
	r := Ap2{
		F: S1{X0: C2{X0: Eq1{X0: Int{V: 0}}, X1: Int{V: 1}}},
		A: Ap2{
			F: B1{X0: Mul1{X0: Int{V: 2}}},
			A: Ap2{
				F: B1{X0: Pwr2{}},
				A: Add1{X0: Int{V: -1}},
			},
		},
	}.Eval(c)
	// log.Printf("%s => %s", t, r)
	return r
}

func (t Pwr2) String() string {
	return "pwr2"
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
