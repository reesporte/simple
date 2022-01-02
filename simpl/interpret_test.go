package simpl

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCreatesInterpreter(t *testing.T) {
	input := `i = 0 
i = i + 1 
if i % 3 == 0 print "fizz" 
if i % 5 == 0 print "buzz" 
if i < 100 goto 2
`
	l := Lexer{}
	l.In = strings.NewReader(input)
	p := Parser{}
	tkns, _ := l.Lex()
	p.Tokens = tkns
	p.Parse()
	i := NewInterpreter(&p.Lines, os.Stdout)
	if !reflect.DeepEqual(*i.Lines, p.Lines) {
		t.Fatalf("expected lines to be equal")
	}
}

func TestOperators(t *testing.T) {
	cases := []struct {
		class    string
		input    string
		expected interface{}
	}{
		{
			class:    "arithmetic",
			input:    "1 + 1",
			expected: 2.0,
		},
		{
			class:    "arithmetic",
			input:    "1 + 1 + 1",
			expected: 3.0,
		},
		{
			class:    "arithmetic",
			input:    "69 * 4 + 5",
			expected: 281.0,
		},
		{
			class:    "arithmetic",
			input:    "69 / 4 - 5",
			expected: 12.25,
		},
		{
			class:    "arithmetic",
			input:    "420 + 69 * 6969 / 3000.4321",
			expected: 580.2639166538713,
		},
		{
			class:    "arithmetic",
			input:    "0.420 + 0.69",
			expected: 1.1099999999999999,
		},
		{
			class:    "arithmetic",
			input:    ".42 + .6",
			expected: 1.02,
		},
		{
			class:    "arithmetic",
			input:    "10 % 3",
			expected: 1,
		},
		{
			class:    "arithmetic",
			input:    "3 % 3",
			expected: 0,
		},
		{
			class:    "arithmetic",
			input:    "420 % 69",
			expected: 6,
		},
		{
			class:    "arithmetic",
			input:    "-420 % 69",
			expected: -6,
		},
		{
			class:    "arithmetic",
			input:    "-420 % -69",
			expected: -6,
		},
		{
			class:    "boolean",
			input:    "42 & 0",
			expected: false,
		},
		{
			class:    "boolean",
			input:    "69 | 0",
			expected: true,
		},
		{
			class:    "boolean",
			input:    "68 & 1 | 5 == 0",
			expected: false,
		},
		{
			class:    "boolean",
			input:    "69 > 2",
			expected: true,
		},
		// bc of the way i parse things you can do weird
		// comparisons of numbers to booleans
		{
			class:    "boolean",
			input:    "-2 > 3 > 5",
			expected: false,
		},
		{
			class:    "boolean",
			input:    "-2 < 3 > 5",
			expected: true,
		},
		{
			class:    "boolean",
			input:    "5 < 3",
			expected: false,
		},
		{
			class:    "boolean",
			input:    "6.0 == 6",
			expected: true,
		},
		{
			class:    "boolean",
			input:    "3 * 2 == 6",
			expected: true,
		},
		{
			class:    "boolean",
			input:    "6 == 3 * 2",
			expected: true,
		},
		{
			class:    "boolean",
			input:    "6 == 4 * 2",
			expected: false,
		},
		{
			class:    "boolean",
			input:    "3 % 3 & 3 % 5",
			expected: false,
		},
		{
			class:    "boolean",
			input:    "1 != 2",
			expected: true,
		},
		{
			class:    "boolean",
			input:    "0 != 2 != 0",
			expected: true,
		},
	}
	l := Lexer{}
	for _, test := range cases {
		t.Run(test.input, func(t *testing.T) {
			l.In = strings.NewReader(test.input)
			tkns, err := l.Lex()
			if err != nil {
				t.Errorf("error '%v' lexing tokens: %v", err, test.input)
			}
			p := Parser{Tokens: tkns}
			p.Parse()
			i := Interpreter{Lines: &p.Lines}
			res := i.Interpret()
			if res != test.expected {
				t.Errorf("expected %v got %v", test.expected, res)
			}
		})
	}
}
