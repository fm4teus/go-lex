package main

import (
	"bufio"
	"fmt"
	"go-lex/lex"
	"io/ioutil"
	"os"
)

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := lex.NewLexer(file)
	b, err := ioutil.ReadAll(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}

	for {
		pos, tok, lit := lexer.Lex(b)
		if tok == lex.EOF {
			break
		}
		fmt.Printf("%d:%d\t%s\t%s\n", pos.Line, pos.Column, tok, lit)
		b = b[pos.Column:]
	}
}
