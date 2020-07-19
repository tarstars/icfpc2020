package interpreter

import (
	"bufio"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode"
)

func makeTokenMap(ts ...Token) map[string]Token {
	r := make(map[string]Token)
	for _, t := range ts {
		r[t.String()] = t
	}
	return r
}

func ParseVarN(s string) Token {
	if !strings.HasPrefix(s, ":") {
		return nil
	}
	n, err := strconv.ParseInt(s[1:], 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return VarN{N: int(n)}
}

func ParseInt(s string) Token {
	idx := strings.IndexFunc(s, unicode.IsDigit)
	sign := strings.HasPrefix(s, "-")
	if (sign && idx != 1) || (!sign && idx != 0) {
		return nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Panic(err)
	}
	return Int{V: n}
}

func IsComment(s string) bool {
	if strings.HasPrefix(s, "#") {
		return true
	}
	return false
}

func ProcessTokens(c Context, toks []string) Token {
	toks = splitOn(toks, "(")
	toks = splitOn(toks, ")")
	toks = splitOn(toks, ",")

	var p Program
	empty := true
	for _, ts := range toks {
		if ts == "(" || ts == "," {
			p.PushProgram(Program{
				Ap{}, Ap{}, Cons{},
			})
			empty = true
			continue
		}
		if ts == ")" {
			if empty {
				p.Push(Nil{})
			}
			p.Push(Nil{})
			continue
		}

		empty = false
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

	// log.Printf("Program: %#v", p)

	tok, err := Interpret(c, p)
	if err != nil {
		log.Panic(err)
	}

	return tok
}

func splitOn(toks []string, sep string) []string {
	var r []string
	for _, t := range toks {
		ts := strings.Split(t, sep)
		if ts[0] != "" {
			r = append(r, ts[0])
		}
		for _, ti := range ts[1:] {
			r = append(r, sep)
			if ti != "" {
				r = append(r, ti)
			}
		}
	}
	return r
}

func ParseLine(c Context, s string) []Token {
	toks := strings.Fields(s)
	if len(toks) == 0 || IsComment(toks[0]) {
		return nil
	}

	assign := false
	var varN VarN

	if len(toks) > 2 && toks[1] == "=" && ParseVarN(toks[0]) != nil {
		varN, assign = ParseVarN(toks[0]).(VarN), true
		toks = toks[2:]
		// log.Printf("Assign %s = %s", varN, strings.Join(toks, " "))
	}

	tok := ProcessTokens(c, toks)

	if assign {
		c.SetVar(varN.N, tok)
		return nil
	}

	r := TailEval(c, tok)
	// log.Printf("%s => %s", tok, r)
	return []Token{r}
}

func ParseReader(c Context, rd io.Reader) []Token {
	lrd := bufio.NewReader(rd)
	var rs []Token
	for line, err := lrd.ReadString('\n'); err != io.EOF || len(line) > 0; line, err = lrd.ReadString('\n') {
		rs = append(rs, ParseLine(c, line)...)
	}
	return rs
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
	Vec{},
	Draw{},
	Checkerboard{},
	Multipledraw{},
	If0{},
	Interact{},
)
