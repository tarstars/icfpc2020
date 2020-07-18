package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
)

type Token interface{}

type Value interface {
	Value() int64
}

type Func interface {
	Apply(v Token) (Program, error)
}

type Ap struct{}

func (t Ap) String() string {
	return "ap"
}

type Var struct {
	N int
}

func (t Var) String() string {
	return fmt.Sprintf(":%d", t.N)
}

type Int struct {
	V int64
}

func (t Int) Value() int64 {
	return t.V
}

func (t Int) String() string {
	return strconv.FormatInt(t.V, 10)
}

func ParseInt(s string) (Int, error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return Int{}, err
	}
	return Int{V: v}, nil
}

type Inc struct{}

func (t Inc) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Int{V: w.Value() + 1}}, nil
}

func (t Inc) String() string {
	return "inc"
}

type Dec struct{}

func (t Dec) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Int{V: w.Value() - 1}}, nil
}

func (t Dec) String() string {
	return "dec"
}

type Add struct{}

func (t Add) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Add1{V: w}}, nil
}

func (t Add) String() string {
	return "add"
}

type Add1 struct {
	V Value
}

func (t Add1) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Int{V: t.V.Value() + w.Value()}}, nil
}

func (t Add1) String() string {
	return fmt.Sprintf("add[%s]", t.V)
}

type Mul struct{}

func (t Mul) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Mul1{V: w}}, nil
}

func (t Mul) String() string {
	return "mul"
}

type Mul1 struct {
	V Value
}

func (t Mul1) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Int{V: t.V.Value() * w.Value()}}, nil
}

func (t Mul1) String() string {
	return fmt.Sprintf("mul[%s]", t.V)
}

type Div struct{}

func (t Div) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	r := Program{Div1{V: w}}
	return r, nil
}

func (t Div) String() string {
	return "div"
}

type Div1 struct {
	V Value
}

func (t Div1) Apply(v Token) (Program, error) {
	w, ok := v.(Value)
	if !ok {
		return nil, fmt.Errorf("Invalid argument: %s(%s)", t, v)
	}
	return Program{Int{V: t.V.Value() / w.Value()}}, nil
}

func (t Div1) String() string {
	return fmt.Sprintf("div[%s]", t.V)
}

type Eq struct{}

func (t Eq) Apply(v Token) (Program, error) {
	return Program{Eq1{X: v}}, nil
}

func (t Eq) String() string {
	return "eq"
}

type Eq1 struct {
	X Token
}

func (t Eq1) Apply(v Token) (Program, error) {
	if reflect.DeepEqual(t.X, v) {
		return Program{True{}}, nil
	}
	return Program{False{}}, nil
}

func (t Eq1) String() string {
	return fmt.Sprintf("eq[%s]", t.X)
}

type Lt struct{}

func (t Lt) Apply(v Token) (Program, error) {
	w := v.(Value)
	return Program{Lt1{X: w}}, nil
}

func (t Lt) String() string {
	return "lt"
}

type Lt1 struct {
	X Value
}

func (t Lt1) Apply(v Token) (Program, error) {
	w := v.(Value)
	if t.X.Value() < w.Value() {
		return Program{True{}}, nil
	}
	return Program{False{}}, nil
}

func (t Lt1) String() string {
	return fmt.Sprintf("lt[%s]", t.X)
}

type True struct{}

func (t True) Apply(v Token) (Program, error) {
	return Program{True1{X: v}}, nil
}

func (t True) String() string {
	return "t"
}

type True1 struct {
	X Token
}

func (t True1) Apply(v Token) (Program, error) {
	return Program{t.X}, nil
}

func (t True1) String() string {
	return fmt.Sprintf("t[%s]", t.X)
}

type False struct{}

func (t False) Apply(v Token) (Program, error) {
	return Program{False1{X: v}}, nil
}

func (t False) String() string {
	return "f"
}

type False1 struct {
	X Token
}

func (t False1) Apply(v Token) (Program, error) {
	return Program{v}, nil
}

func (t False1) String() string {
	return fmt.Sprintf("f[%s]", t.X)
}
