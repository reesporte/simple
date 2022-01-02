package simpl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
)

// TokenType represents the type of token found
type TokenType uint64

// Various token types
const (
	Operator TokenType = iota
	Str
	Num
	Assignment
	Boolop
	Builtin
	Keyword
	Var
	Paren
	Newline
)

func (t TokenType) String() string {
	return [...]string{"operator", "str", "num", "assignment", "boolop", "builtin", "keyword", "variable", "parenthesis", "newline"}[t]
}

// Token holds information about a token
type Token struct {
	Class TokenType
	Repr  string
}

// Lexer holds state needed for lexing
type Lexer struct {
	In io.Reader
}

// Lex returns the lexed tokens in an io.Reader
func (l *Lexer) Lex() (tkns []Token, errors []error) {
	input := bufio.NewReader(l.In)
	tkn := []rune{}
	quotes := 0
	escape := false
	comment := false
	for {
		if c, _, err := input.ReadRune(); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		} else {
			if escape && !comment {
				switch c {
				case 'n':
					tkn = append(tkn, '\n')
				case 't':
					tkn = append(tkn, '\t')
				case '"':
					tkn = append(tkn, '"')
				case '\'':
					tkn = append(tkn, '\'')
				case '\\':
					tkn = append(tkn, '\\')
				default:
					errors = append(errors, fmt.Errorf("unknown escape: %v", c))
				}
				escape = false
				continue
			} else {
				switch c {
				case '\\':
					// found an escape character
					escape = true
				case ' ', '\n':
					if quotes == 0 {
						if len(tkn) > 0 && !comment {
							strTkn := string(tkn)
							tkn = []rune{}
							class, err := classifyToken(strTkn)
							if err != nil {
								errors = append(errors, err)
								continue
							}
							tkns = append(tkns, Token{Class: class, Repr: strTkn})
						}
						if c == '\n' {
							comment = false
							tkns = append(tkns, Token{Class: Newline, Repr: "\\n"})
						}
					} else if !comment {
						tkn = append(tkn, c)
					}
				case '"':
					quotes++
					if quotes == 2 && !comment {
						strTkn := string(tkn)
						tkn = []rune{}
						if err != nil {
							errors = append(errors, err)
							continue
						}
						tkns = append(tkns, Token{Class: Str, Repr: strTkn})
						quotes = 0
					}
				case '#':
					comment = true
				default:
					if !comment {
						tkn = append(tkn, c)
					}
				}
			}
		}
	}
	if len(tkn) > 0 && !comment {
		strTkn := string(tkn)
		class, err := classifyToken(strTkn)
		if err != nil {
			errors = append(errors, err)
		}
		tkns = append(tkns, Token{Class: class, Repr: strTkn})
	}

	return tkns, errors
}

func classifyToken(t string) (TokenType, error) {
	switch t {
	case "+", "-", "*", "/", "%":
		return Operator, nil
	case "<", ">", "==", "&", "|", "!=":
		return Boolop, nil
	case "print", "goto":
		return Builtin, nil
	case "if":
		return Keyword, nil
	case "=":
		return Assignment, nil
	case "(", ")":
		return Paren, nil
	}
	if t[0] == '"' && t[len(t)-1] == '"' {
		return Str, nil
	}
	_, err := strconv.ParseFloat(t, 64)
	if err == nil {
		return Num, nil
	}
	// if it's not anything else, it's probably an ident
	if isAlphaNumeric(t) {
		return Var, nil
	}
	return 0, fmt.Errorf("unrecognized token: '%v'", t)
}

func isAlphaNumeric(tkn string) bool {
	for _, r := range tkn {
		if !(('0' <= r && r <= '9') || ('A' <= r && r <= 'z')) {
			return false
		}
	}
	return true
}
