package main

import (
	"bufio"
	"fmt"
	"go-lex/lex"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var wsRegex = regexp.MustCompile(`\S`)
var newLineRegex = regexp.MustCompile(`\n`)

func calculateLineColumn(newLineIndexes []int, index int) (line, column int) {
	lastIdx := -1
	for _, i := range newLineIndexes {
		if i > index {
			column = index - lastIdx
			break
		}
		lastIdx = i
		line++
	}
	return line + 1, column
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func trimLeftSpace(s string) string {
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			return strings.TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	stop := len(s)

	return s[start:stop]
}

func main() {
	args := os.Args
	var filename string
	if len(args) > 1 {
		filename = args[1]
	}
	if filename == "" {
		panic(fmt.Errorf("invalid file name"))
	}

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	lexer := lex.NewLexer(file)
	b, err := ioutil.ReadAll(bufio.NewReader(file))
	if err != nil {
		panic(err)
	}

	var index int
	var newLineIndexes []int
	errs := []struct {
		line int
		col  int
	}{}

	newLineChars := newLineRegex.FindAllIndex(b, -1)
	for _, brIndex := range newLineChars {
		newLineIndexes = append(newLineIndexes, brIndex[0])
	}

	for {
		p := wsRegex.FindIndex(b)
		if p != nil {
			index += p[0]
		}

		b = []byte(trimLeftSpace(string(b)))

		end, tok, lit := lexer.Lex(b)
		if tok == lex.EOF {
			break
		}

		l, c := calculateLineColumn(newLineIndexes, index)

		if tok == lex.ERROR {
			errs = append(errs, struct {
				line int
				col  int
			}{l, c})
		}

		fmt.Printf("%d:%d\t%s\t%s\n", l, c, tok, lit)
		index += end
		b = b[end:]
	}

	if len(errs) > 0 {
		fmt.Printf("\nfound %d errors: \n", len(errs))
		for _, e := range errs {
			fmt.Printf("line: %d \t col: %d\n", e.line, e.col)
		}
	}
}
