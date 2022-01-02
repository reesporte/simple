package simpl

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestIsAlphaNumeric(t *testing.T) {
	tests := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuwxyz0123456789"
	for _, v := range tests {
		s := string(v)
		if !isAlphaNumeric(s) {
			t.Fatalf("expected %s to be alphanumeric", s)
		}
	}
}

func TestLex(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []Token
		errors   []error
	}{
		{
			name: "fizzbuzz",
			input: `i = 0 
i = i + 1 
if i % 3 == 0 print "fizz" 
if i % 5 == 0 print "buzz" 
if i % 3 & i % 5 print i
print "\n"
if i < 100 goto 2 
`,
			expected: []Token{
				Token{Class: Var, Repr: "i"}, Token{Class: Assignment, Repr: "="}, Token{Class: Num, Repr: "0"}, Token{Class: Newline, Repr: "\\n"},
				Token{Class: Var, Repr: "i"}, Token{Class: Assignment, Repr: "="}, Token{Class: Var, Repr: "i"}, Token{Class: Operator, Repr: "+"},
				Token{Class: Num, Repr: "1"}, Token{Class: Newline, Repr: "\\n"}, Token{Class: Keyword, Repr: "if"}, Token{Class: Var, Repr: "i"},
				Token{Class: Operator, Repr: "%"}, Token{Class: Num, Repr: "3"}, Token{Class: Boolop, Repr: "=="}, Token{Class: Num, Repr: "0"},
				Token{Class: Builtin, Repr: "print"}, Token{Class: Str, Repr: "fizz"}, Token{Class: Newline, Repr: "\\n"}, Token{Class: Keyword, Repr: "if"},
				Token{Class: Var, Repr: "i"}, Token{Class: Operator, Repr: "%"}, Token{Class: Num, Repr: "5"}, Token{Class: Boolop, Repr: "=="},
				Token{Class: Num, Repr: "0"}, Token{Class: Builtin, Repr: "print"}, Token{Class: Str, Repr: "buzz"}, Token{Class: Newline, Repr: "\\n"},
				Token{Class: Keyword, Repr: "if"}, Token{Class: Var, Repr: "i"}, Token{Class: Operator, Repr: "%"}, Token{Class: Num, Repr: "3"},
				Token{Class: Boolop, Repr: "&"}, Token{Class: Var, Repr: "i"}, Token{Class: Operator, Repr: "%"}, Token{Class: Num, Repr: "5"},
				Token{Class: Builtin, Repr: "print"}, Token{Class: Var, Repr: "i"}, Token{Class: Newline, Repr: "\\n"}, Token{Class: Builtin, Repr: "print"},
				Token{Class: Str, Repr: "\n"}, Token{Class: Newline, Repr: "\\n"}, Token{Class: Keyword, Repr: "if"}, Token{Class: Var, Repr: "i"},
				Token{Class: Boolop, Repr: "<"}, Token{Class: Num, Repr: "100"}, Token{Class: Builtin, Repr: "goto"}, Token{Class: Num, Repr: "2"},
				Token{Class: Newline, Repr: "\\n"}},
			errors: []error{},
		},
		{
			name:  "not equals",
			input: "i != 0",
			expected: []Token{
				Token{Class: Var, Repr: "i"},
				Token{Class: Boolop, Repr: "!="},
				Token{Class: Num, Repr: "0"},
			},
		},
		{
			name:  "parenthesis",
			input: "( ) ( )",
			expected: []Token{
				Token{Class: Paren, Repr: "("},
				Token{Class: Paren, Repr: ")"},
				Token{Class: Paren, Repr: "("},
				Token{Class: Paren, Repr: ")"},
			},
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			l := Lexer{In: strings.NewReader(test.input)}
			results, errs := l.Lex()
			if len(results) != len(test.expected) {
				t.Errorf("wrong number of results: expected %v, got %v.\nexpected contents: %v,\ngot      contents: %v", len(test.expected), len(results), test.expected, results)
			}
			if len(errs) != len(test.errors) {
				t.Errorf("wrong number of errors: expected %v, got %v", len(test.errors), len(errs))
			}
			for i, e := range errs {
				if !reflect.DeepEqual(e, test.errors[i]) {
					t.Errorf("expected %v, got %v", test.errors[i], e)
				}
			}
			for i, r := range results {
				if !reflect.DeepEqual(r, test.expected[i]) {
					t.Errorf("expected %v, got %v", test.expected[i], r)
				}
			}
		})
	}
}

func TestClassifyToken(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
		err      error
	}{
		{
			input:    "<",
			expected: Boolop,
			err:      nil,
		},
		{
			input:    ">",
			expected: Boolop,
			err:      nil,
		},
		{
			input:    "+",
			expected: Operator,
			err:      nil,
		},
		{
			input:    "-",
			expected: Operator,
			err:      nil,
		},
		{
			input:    "*",
			expected: Operator,
			err:      nil,
		},
		{
			input:    "/",
			expected: Operator,
			err:      nil,
		},
		{
			input:    "%",
			expected: Operator,
			err:      nil,
		},
		{
			input:    "==",
			expected: Boolop,
			err:      nil,
		},
		{
			input:    "&",
			expected: Boolop,
			err:      nil,
		},
		{
			input:    "|",
			expected: Boolop,
			err:      nil,
		},
		{
			input:    "##|",
			expected: 0,
			err:      fmt.Errorf("unrecognized token: '%v'", "##|"),
		},
		{
			input:    "print",
			expected: Builtin,
			err:      nil,
		},
		{
			input:    "goto",
			expected: Builtin,
			err:      nil,
		},
		{
			input:    "if",
			expected: Keyword,
			err:      nil,
		},
		{
			input:    "=",
			expected: Assignment,
			err:      nil,
		},
		{
			input:    "\"##|\"",
			expected: Str,
			err:      nil,
		},
		{
			input:    "-234567890",
			expected: Num,
			err:      nil,
		},
		{
			input:    "67890198763431",
			expected: Num,
			err:      nil,
		},
		{
			input:    "420695",
			expected: Num,
			err:      nil,
		},
		{
			input:    "v",
			expected: Var,
			err:      nil,
		},
		{
			input:    "variablenmae",
			expected: Var,
			err:      nil,
		},
		{
			input:    "asdf-=asdf",
			expected: 0,
			err:      fmt.Errorf("unrecognized token: '%v'", "asdf-=asdf"),
		},
		{
			input:    "!=",
			expected: Boolop,
			err:      nil,
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			class, err := classifyToken(test.input)
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("expected  %v, got %v", test.err, err)
			}
			if class != test.expected {
				t.Errorf("expected %v, got %v", test.expected, class)
			}
		})
	}

}
