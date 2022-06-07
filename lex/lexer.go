package lex

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"regexp"
)

type Token int

var (
	separator       = `[:;,\s\(\)\[\]\{\}]`
	operator        = `[+\-*\/<>!=]`
	separatorRegex  = regexp.MustCompile(fmt.Sprintf(`^%s`, separator))
	stringRegex     = regexp.MustCompile(fmt.Sprintf(`^\s*"[^"]*"(%s|%s)`, separator, operator))
	numberRegex     = regexp.MustCompile(fmt.Sprintf(`^\s*[0-9]+(\.[0-9]+)?(%s|%s)`, separator, operator))
	identifierRegex = regexp.MustCompile(fmt.Sprintf(`^[A-Za-z_][A-Za-z0-9_]*(%s|%s)`, separator, operator))
	operatorRegex   = regexp.MustCompile(fmt.Sprintf(`^((%s=?)|(\+\+)|(--)|(\|\|)|(&&))(%s|[0-9A-Za-z"])`, operator, separator))
)

const (
	EOF = iota
	ERROR
	NUM
	STRING
	SEP
	IDENTIFIER
	KEYWORD
	OPERATOR
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	ASSIGN
	OP_ADIT
	OP_MULT
	OP_LOGIC
	KEYWORD_IF
	KEYWORD_ELSE
	SEMI
	KEYWORD_FOR
	TYPE
)

var keywords = []string{
	"for",
	"if",
	"switch",
	"case",
	"else",
	"return",
	"int",
	"float",
	"char",
}

var tokens = []string{
	EOF:          "EOF",
	ERROR:        "ERROR\t",
	NUM:          "NUM\t",
	STRING:       "STRING\t",
	SEP:          "SEPARATOR",
	IDENTIFIER:   "IDENTIFIER",
	KEYWORD:      "KEYWORD\t",
	OPERATOR:     "OPERATOR",
	LPAREN:       "LPAREN\t",
	RPAREN:       "RPAREN\t",
	LBRACE:       "LBRACE\t",
	RBRACE:       "RBRACE\t",
	ASSIGN:       "ASSIGN\t",
	OP_ADIT:      "OP_ADIT",
	OP_MULT:      "OP_MULT",
	OP_LOGIC:     "OP_LOGIC",
	KEYWORD_IF:   "KEYWORD_IF",
	KEYWORD_ELSE: "KEYWORD_ELSE",
	SEMI:         "SEMI",
	KEYWORD_FOR:  "KEYWORD_FOR",
	TYPE:         "TYPE",
}

type responseToken struct {
	End   int
	Token Token
	Lit   string
}

func (t Token) String() string {
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

	ch := make(chan responseToken)

	wg := &sync.WaitGroup{}
	wg.Add(5)
	go matchRegexp(stringRegex, b, STRING, ch, wg)
	go matchRegexp(numberRegex, b, NUM, ch, wg)
	go matchRegexp(separatorRegex, b, SEP, ch, wg)
	go matchRegexp(identifierRegex, b, IDENTIFIER, ch, wg)
	go matchRegexp(operatorRegex, b, OPERATOR, ch, wg)
	go func() {
		wg.Wait()
		close(ch)
	}()

	res, ok := <-ch
	if !ok {
		return 1, ERROR, "Error"
	}
	return res.End, res.Token, res.Lit
}

func matchRegexp(re *regexp.Regexp, b []byte, token Token, ch chan responseToken, wg *sync.WaitGroup) {
	if len(strings.TrimSpace(string(b))) == 0 {
		ch <- responseToken{0, EOF, tokens[EOF]}
		return
	}

	r := re.FindIndex(b)
	var notSeparator int
	if token != SEP {
		notSeparator = 1
	}

	if r != nil {
		var lit = string(b[r[0] : r[1]-notSeparator])

		switch token {
		case IDENTIFIER:
			for _, k := range keywords {
				if lit == k {
					token = KEYWORD
					switch lit {
					case "if":
						token = KEYWORD_IF
					case "else":
						token = KEYWORD_ELSE
					case "for":
						token = KEYWORD_FOR
					case "int", "float", "char":
						token = TYPE
					}
					ch <- responseToken{End: r[1] - notSeparator, Token: token, Lit: lit}
					return
				}
			}
		case SEP:
			switch lit {
			case "(":
				token = LPAREN
			case ")":
				token = RPAREN
			case "{":
				token = LBRACE
			case "}":
				token = RBRACE
			case ";":
				token = SEMI
			}
		case OPERATOR:
			switch lit {
			case "=":
				token = ASSIGN
			case "+", "-":
				token = OP_ADIT
			case "*", "/":
				token = OP_MULT
			case "==", "!=", "!", ">", "<", "<=", ">=":
				token = OP_LOGIC
			}
		}

		ch <- responseToken{End: r[1] - notSeparator, Token: token, Lit: lit}
		return
	}

	wg.Done()
}
