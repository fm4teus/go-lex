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

// TrimSpace returns a slice of the string s, with all leading
// and trailing white space removed, as defined by Unicode.
func trimLeftSpace(s string) string {
	// Fast path for ASCII: look for the first ASCII non-space byte
	start := 0
	for ; start < len(s); start++ {
		c := s[start]
		if c >= utf8.RuneSelf {
			// If we run into a non-ASCII byte, fall back to the
			// slower unicode-aware method on the remaining bytes
			return strings.TrimFunc(s[start:], unicode.IsSpace)
		}
		if asciiSpace[c] == 0 {
			break
		}
	}

	// Now look for the first ASCII non-space byte from the end
	stop := len(s)

	// At this point s[start:stop] starts and ends with an ASCII
	// non-space bytes, so we're done. Non-ASCII cases have already
	// been handled above.
	return s[start:stop]
}

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

	var index int
	var newLineIndexes []int

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

		fmt.Printf("%d:%d\t%s\t%s\n", l, c, tok, lit)
		// b = []byte(strings.TrimSpace(string(b[end:])))
		index += end
		b = b[end:]
	}
}
