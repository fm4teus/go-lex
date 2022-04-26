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
	operatorRegex   = regexp.MustCompile(fmt.Sprintf(`^((%s=?)|(\+\+)|(--))(%s|[0-9A-Za-z"])`, operator, separator))
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
)

var m sync.Mutex

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
	EOF:        "EOF",
	ERROR:      "ERROR\t",
	NUM:        "NUM\t",
	STRING:     "STRING\t",
	SEP:        "SEPARATOR",
	IDENTIFIER: "IDENTIFIER",
	KEYWORD:    "KEYWORD\t",
	OPERATOR:   "OPERATOR",
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
		if token == IDENTIFIER {
			for _, k := range keywords {
				if lit == k {
					ch <- responseToken{End: r[1] - notSeparator, Token: KEYWORD, Lit: lit}
					return
				}
			}
		}
		ch <- responseToken{End: r[1] - notSeparator, Token: token, Lit: lit}
		return
	}

	m.Lock()
	wg.Done()
	m.Unlock()
}
