package lex

import (
	"fmt"
	"io"
	"strings"

	"regexp"
)

type Token int

var (
	separator       = `[;,\s\(\)\[\]\{\}]`
	separatorRegex  = regexp.MustCompile(`^` + separator)
	stringRegex     = regexp.MustCompile(`^\s*".*"` + separator)
	numberRegex     = regexp.MustCompile(`^\s*[0-9]+(\.[0-9]+)?` + separator)
	identifierRegex = regexp.MustCompile(`^[A-Za-z_]\w*` + separator)
)

const (
	EOF = iota
	ERROR
	NUM
	STRING
	SEP
	IDENTIFIER

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =
)

var tokens = []string{
	EOF:        "EOF",
	ERROR:      "ERROR",
	NUM:        "NUM",
	STRING:     "STRING",
	SEP:        "SEPARATOR",
	IDENTIFIER: "IDENTIFIER",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	ASSIGN: "=",
}

func (t Token) String() string {
	if t == -1 {
		return "ERROR"
	}
	return tokens[t]
}

type Position struct {
	Line   int
	Column int
}

type Lexer struct{}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{}
}

func (l *Lexer) Lex(b []byte) (int, Token, string) {

	p := stringRegex.FindIndex(b)

	if p != nil {
		return p[1] - 1, STRING, string(b[p[0] : p[1]-1])
	}

	q := numberRegex.FindIndex(b)

	if q != nil {
		return q[1] - 1, NUM, string(b[q[0] : q[1]-1])
	}

	r := separatorRegex.FindIndex(b)
	if r != nil {
		return r[1], SEP, string(b[r[0] : r[1]-1])
	}

	s := identifierRegex.FindIndex(b)
	if s != nil {
		return s[1] - 1, IDENTIFIER, string(b[s[0] : s[1]-1])
	}

	if len(strings.TrimSpace(string(b))) == 0 {
		return 0, EOF, tokens[EOF]
	}
	fmt.Println("Err: ", string(b))
	return 1, -1, "Error"
}
