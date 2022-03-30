package main

import (
	"fmt"
	"os"

	"go-lex/lex"
)

func main() {
	file, err := os.Open("input.test")
	if err != nil {
		panic(err)
	}

	lexer := lex.NewLexer(file)

	for {
		pos, tok, lit := lexer.Lex()
		if tok == lex.EOF {
			break
		}
		fmt.Printf("%d:%d\t%s\t%s\n", pos.Line, pos.Column, tok, lit)
	}
}
