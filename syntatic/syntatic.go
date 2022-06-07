package syntatic

import (
	"errors"
	"fmt"
	"go-lex/lex"
)

// math
// E -> TE'
// E' -> +TE'
// E' -> vazio
// T-> FT'
// T' -> *FT'
// T' -> vazio
// F -> (E)
// F-> número||id

type Syntatic struct {
	tokenTable []*TokenInfo
	index      uint
	debug      bool
}

type TokenInfo struct {
	Line  int
	Col   int
	Token lex.Token
}

func NewSyntatic(tokenTable []*TokenInfo, debug bool) *Syntatic {
	return &Syntatic{
		tokenTable: tokenTable,
		index:      0,
		debug:      debug,
	}
}

func (s *Syntatic) presentToken() TokenInfo {
	if len(s.tokenTable) > int(s.index) {
		return *s.tokenTable[s.index]
	}
	return TokenInfo{}
}

func (s *Syntatic) ReadToken() error {

	if int(s.index) > len(s.tokenTable)-1 {
		return errors.New("EOF")
	}
	if s.debug {
		fmt.Printf("READ: %s\n", s.presentToken().Token)
	}
	s.index++

	return nil
}

func (s *Syntatic) Analyse() error {
	if len(s.tokenTable) == 0 {
		return nil
	}

	if err := s.S(); err != nil {
		return err
	}

	if len(s.tokenTable) > int(s.index) {
		return fmt.Errorf("failed to read token %v", s.presentToken())
	}

	return nil
}

// S -> A S
// S -> I S
// S -> FL S
// S -> vazio
func (s *Syntatic) S() error {
	switch s.presentToken().Token {
	case lex.IDENTIFIER, lex.TYPE:
		if err := s.A(); err != nil {
			return err
		}
		if s.presentToken().Token != lex.SEMI {
			return fmt.Errorf("expected ';'; found %v instead", s.presentToken())
		}
		if err := s.ReadToken(); err != nil {
			return nil
		}
	case lex.KEYWORD_IF:
		if err := s.I(); err != nil {
			return err
		}
	case lex.KEYWORD_FOR:
		if err := s.FL(); err != nil {
			return err
		}
	default:
		return nil
	}
	if err := s.S(); err != nil {
		return err
	}
	return nil
}

// A  -> A' id = E
// A' -> type
// A' -> vazio
func (s *Syntatic) A() error {
	if s.debug {
		fmt.Println("   ---   begin assign   ---   ")
	}
	s.A_()

	if s.presentToken().Token != lex.IDENTIFIER {
		return fmt.Errorf("expected Identifier; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if s.presentToken().Token != lex.ASSIGN {
		return fmt.Errorf("expected =; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.E(); err != nil {
		return err
	}

	if s.debug {
		fmt.Println("   ---   end assign   ---   ")
	}
	return nil
}

func (s *Syntatic) A_() error {
	if s.presentToken().Token != lex.TYPE {
		return nil
	}
	s.ReadToken()
	return nil
}

func (s *Syntatic) E() error {
	if err := s.T(); err != nil {
		return err
	}

	if err := s.E_(); err != nil {
		return err
	}
	return nil
}

func (s *Syntatic) T() error {
	if err := s.F(); err != nil {
		return err
	}

	if err := s.T_(); err != nil {
		return err
	}

	return nil
}

func (s *Syntatic) E_() error {
	if s.presentToken().Token != lex.OP_ADIT {
		return nil
		// or empty
		// return fmt.Errorf("not found token +")
	}
	// read token + or -
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.T(); err != nil {
		return err
	}

	if err := s.E_(); err != nil {
		return err
	}

	return nil
}

func (s *Syntatic) T_() error {
	if s.presentToken().Token != lex.OP_MULT {
		// or empty ?
		return nil
		// return fmt.Errorf("not found token * at line: %d, col: %d", t.Line, t.Col)
	}
	// read token * or /
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.F(); err != nil {
		return err
	}

	if err := s.T_(); err != nil {
		return err
	}

	return nil
}

func (s *Syntatic) F() error {
	if s.presentToken().Token == lex.NUM || s.presentToken().Token == lex.IDENTIFIER {
		// read token
		s.ReadToken()
		return nil
	}

	if s.presentToken().Token != lex.LPAREN {
		return fmt.Errorf("not found token ( at line: %d, col: %d", s.presentToken().Line, s.presentToken().Col)
	}
	// read token (
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.E(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.RPAREN {
		return fmt.Errorf("not found token ); found %v instead", s.presentToken())
	}
	// read token )
	s.ReadToken()

	return nil
}

// logic TLDR
// L  -> M
// L  -> ML'
// L' -> [||&&<><=>=]ML'
// L' -> vazio
// M  -> NM'
// M' -> &&NM'
// M' -> vazio
// N  -> PN'
// N' -> ==PN'
// N' -> vazio
// P  -> QP'
// P' -> <=QP'
// P' -> vazio
// Q  -> !R
// Q  -> R
// R  -> (L)
// R  -> id/num

// logic
// L  -> ML'
// L' -> logicML'
// L' -> vazio
// M  -> (L)
// M  -> id/num

func (s *Syntatic) L() error {
	if s.debug {
		fmt.Println("   ---   begin logic   ---   ")
	}
	if err := s.M(); err != nil {
		return err
	}
	if err := s.L_(); err != nil {
		return err
	}
	if s.debug {
		fmt.Println("   ---   end logic   ---   ")
	}
	return nil
}

func (s *Syntatic) M() error {
	if s.presentToken().Token == lex.IDENTIFIER || s.presentToken().Token == lex.NUM {
		// read token
		s.ReadToken()
		return nil
	}
	if s.presentToken().Token != lex.LPAREN {
		return fmt.Errorf("not found token (; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if s.presentToken().Token != lex.RPAREN {
		return fmt.Errorf("not found token ); found %v instead", s.presentToken())
	}

	s.ReadToken()

	return nil
}

func (s *Syntatic) L_() error {
	if s.presentToken().Token != lex.OP_LOGIC {
		// or empty
		return nil
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.M(); err != nil {
		return err
	}
	if err := s.L_(); err != nil {
		return err
	}
	return nil
}

// if-else
// I  -> if(L){S}J
// J  -> elseJ'
// J  -> vazio
// J' -> I
// J' -> {S}

func (s *Syntatic) I() error {
	if s.debug {
		fmt.Println("   ---   begin if   ---   ")
	}
	if s.presentToken().Token != lex.KEYWORD_IF {
		return fmt.Errorf("not found token IF; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if s.presentToken().Token != lex.LPAREN {
		return fmt.Errorf("not found token (; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.L(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.RPAREN {
		return fmt.Errorf("not found token ); found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if s.presentToken().Token != lex.LBRACE {
		return fmt.Errorf("not found token {; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.S(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.RBRACE {
		return fmt.Errorf("not found token }; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.J(); err != nil {
		return err
	}
	if s.debug {
		fmt.Println("   ---   end if   ---   ")
	}
	return nil
}
func (s *Syntatic) J() error {
	if s.presentToken().Token != lex.KEYWORD_ELSE {
		return nil
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.J_(); err != nil {
		return err
	}

	return nil
}
func (s *Syntatic) J_() error {
	if s.presentToken().Token == lex.KEYWORD_IF {
		return s.I()
	}

	if s.presentToken().Token != lex.LBRACE {
		return fmt.Errorf("not found token {; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.S(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.RBRACE {
		return fmt.Errorf("not found token }; found %v instead", s.presentToken())
	}
	s.ReadToken()

	return nil
}

// FL -> for(A;L;A){S}
func (s *Syntatic) FL() error {
	if s.debug {
		fmt.Println("   ---   begin for   ---   ")
	}
	if s.presentToken().Token != lex.KEYWORD_FOR {
		return fmt.Errorf("expected token 'for'; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if s.presentToken().Token != lex.LPAREN {
		return fmt.Errorf("expected token '('; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.A(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.SEMI {
		return fmt.Errorf("expected token ';'; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.L(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.SEMI {
		return fmt.Errorf("expected token ';'; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.A(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.RPAREN {
		return fmt.Errorf("expected token ')'; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if s.presentToken().Token != lex.LBRACE {
		return fmt.Errorf("expected token '{'; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}

	if err := s.S(); err != nil {
		return err
	}

	if s.presentToken().Token != lex.RBRACE {
		return fmt.Errorf("expected token '}'; found %v instead", s.presentToken())
	}
	if err := s.ReadToken(); err != nil {
		return nil
	}
	if s.debug {
		fmt.Println("   ---   end for   ---   ")
	}

	return nil
}
