package lex

import (
	"io"

	"regexp"
)

type Token int

var (
	stringRegex = regexp.MustCompile(`".*"`)
	numberRegex = regexp.MustCompile(`[0-9]+(\.[0-9]+)?`)
)

const (
	EOF = iota
	ILLEGAL
	IDENT
	NUM
	STRING
	SEMI // ;

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	ASSIGN // =
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	NUM:     "NUM",
	STRING:  "STRING",
	SEMI:    ";",

	// Infix ops
	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	ASSIGN: "=",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	Line   int
	Column int
}

type Lexer struct {
	pos Position
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos: Position{Line: 1, Column: 0},
	}
}

func (l *Lexer) Lex(b []byte) (Position, Token, string) {

	p := stringRegex.FindIndex(b)

	if p != nil {
		pos := Position{Line: p[0], Column: p[1]}
		return pos, STRING, string(b[pos.Line:pos.Column])
	}

	q := numberRegex.FindIndex(b)

	if q != nil {
		pos := Position{Line: q[0], Column: q[1]}
		return pos, NUM, string(b[pos.Line:pos.Column])
	}

	return Position{}, EOF, tokens[EOF]

}
