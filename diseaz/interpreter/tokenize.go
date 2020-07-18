package interpreter

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// var spaces = regexp.MustCompile("[ \t]+")

// func splitTokens(s string) []string {
// 	return spaces.Split(strings.TrimSpace(s), -1)
// }

func makeTokenMap(ts ...Token) map[string]Token {
	r := make(map[string]Token)
	for _, t := range ts {
		r[t.String()] = t
	}
	return r
}

func ParseVarN(s string) Token {
	r0, rl := utf8.DecodeRuneInString(s)
	if r0 != ':' {
		return nil
	}
	n, err := strconv.ParseInt(s[rl:], 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return VarN{N: int(n)}
}

func ParseInt(s string) Token {
	r0, _ := utf8.DecodeRuneInString(s)
	if !unicode.IsDigit(r0) && r0 != '-' {
		return nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return Int{V: n}
}

func IsComment(s string) bool {
	r0, _ := utf8.DecodeRuneInString(s)
	if r0 == '#' {
		return true
	}
	return false
}

func ParseLine(c Context, s string) Token {
	toks := strings.Fields(s)
	if len(toks) == 0 || IsComment(toks[0]) {
		return nil
	}

	assign := false
	var varN VarN

	if len(toks) > 2 && toks[1] == "=" && ParseVarN(toks[0]) != nil {
		varN, assign = ParseVarN(toks[0]).(VarN), true
		toks = toks[2:]
	}

	var p Program
	for _, ts := range toks {
		t := ParseVarN(ts)
		if t != nil {
			p.Push(t)
			continue
		}
		t = ParseInt(ts)
		if t != nil {
			p.Push(t)
			continue
		}
		t = tokenMap[ts]
		if t != nil {
			p.Push(t)
			continue
		}
		log.Panicf("Unknown token: %#v", ts)
	}

	tok, err := Interpret(c, p)
	if err != nil {
		log.Panic(err)
	}

	if assign {
		c.SetVar(varN.N, tok)
		return nil
	}

	r := tok.Eval(c)
	// log.Printf("%s => %s", tok, r)
	return r
}

func ParseReader(c Context, rd io.Reader) Token {
	lrd := bufio.NewReader(rd)
	var r Token
	for line, err := lrd.ReadString('\n'); err != io.EOF || len(line) > 0; line, err = lrd.ReadString('\n') {
		r = ParseLine(c, line)
	}
	return r
}

var tokenMap = makeTokenMap(
	Ap{},
	Inc{},
	Dec{},
	Add{},
	Mul{},
	Div{},
	Eq{},
	Lt{},
	Modulate{},
	Demodulate{},
	Send{},
	Neg{},
	S{},
	C{},
	B{},
	Pwr2{},
	I{},
	True{},
	False{},
	Cons{},
	Car{},
	Cdr{},
	Nil{},
	IsNil{},
	Draw{},
	Checkerboard{},
	Multipledraw{},
	If0{},
)
