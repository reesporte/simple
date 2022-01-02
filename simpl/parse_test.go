package simpl

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input    string
		expected []*Node
	}{
		{
			input: "if ( i % 3 != 0 ) & ( i % 5 != 0 ) print i ",
			expected: []*Node{
				&Node{
					left: &Node{
						left: &Node{
							left: &Node{
								left:  &Node{val: Token{Class: Var, Repr: "i"}},
								right: &Node{val: Token{Class: Num, Repr: "3"}},
								val:   Token{Class: Operator, Repr: "%"},
							},
							right: &Node{val: Token{Class: Num, Repr: "0"}},
							val:   Token{Class: Boolop, Repr: "!="},
						},
						right: &Node{
							left: &Node{
								left:  &Node{val: Token{Class: Var, Repr: "i"}},
								right: &Node{val: Token{Class: Num, Repr: "5"}},
								val:   Token{Class: Operator, Repr: "%"},
							},
							right: &Node{val: Token{Class: Num, Repr: "0"}},
							val:   Token{Class: Boolop, Repr: "!="},
						},
						val: Token{Class: Boolop, Repr: "&"},
					},
					right: &Node{
						right: &Node{val: Token{Class: Var, Repr: "i"}},
						val:   Token{Class: Builtin, Repr: "print"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
			},
		},
		{
			input: "i + 0",
			expected: []*Node{
				&Node{
					left: &Node{
						val: Token{Class: Var, Repr: "i"},
					},
					right: &Node{
						val: Token{Class: Num, Repr: "0"},
					},
					val: Token{Class: Operator, Repr: "+"},
				},
			},
		},
		{
			input: "1 & 1 | 2 == 3",
			expected: []*Node{
				&Node{
					left: &Node{
						left: &Node{
							val: Token{Class: Num, Repr: "1"},
						},
						right: &Node{
							left: &Node{
								val: Token{Class: Num, Repr: "1"},
							},
							right: &Node{
								val: Token{Class: Num, Repr: "2"},
							},
							val: Token{Class: Boolop, Repr: "|"},
						},
						val: Token{Class: Boolop, Repr: "&"},
					},
					right: &Node{
						val: Token{Class: Num, Repr: "3"},
					},
					val: Token{Class: Boolop, Repr: "=="},
				},
			},
		},
		{
			input: "i = 0",
			expected: []*Node{
				&Node{
					left: &Node{
						val: Token{Class: Var, Repr: "i"},
					},
					right: &Node{
						val: Token{Class: Num, Repr: "0"},
					},
					val: Token{Class: Assignment, Repr: "="},
				},
			},
		},
		{
			input: "i = i + 1",
			expected: []*Node{
				&Node{
					left: &Node{
						val: Token{Class: Var, Repr: "i"},
					},
					right: &Node{
						left:  &Node{val: Token{Class: Var, Repr: "i"}},
						right: &Node{val: Token{Class: Num, Repr: "1"}},
						val:   Token{Class: Operator, Repr: "+"},
					},
					val: Token{Class: Assignment, Repr: "="},
				},
			},
		},
		{
			input: "i % 3 == 0",
			expected: []*Node{
				&Node{
					left: &Node{
						left:  &Node{val: Token{Class: Var, Repr: "i"}},
						right: &Node{val: Token{Class: Num, Repr: "3"}},
						val:   Token{Class: Operator, Repr: "%"},
					},
					right: &Node{val: Token{Class: Num, Repr: "0"}},
					val:   Token{Class: Boolop, Repr: "=="},
				},
			},
		},
		{
			input: "i % 3 != 0",
			expected: []*Node{
				&Node{
					left: &Node{
						left:  &Node{val: Token{Class: Var, Repr: "i"}},
						right: &Node{val: Token{Class: Num, Repr: "3"}},
						val:   Token{Class: Operator, Repr: "%"},
					},
					right: &Node{val: Token{Class: Num, Repr: "0"}},
					val:   Token{Class: Boolop, Repr: "!="},
				},
			},
		},

		{
			input: "i % 3 & i % 5",
			expected: []*Node{
				&Node{
					left: &Node{
						left:  &Node{val: Token{Class: Var, Repr: "i"}},
						right: &Node{val: Token{Class: Num, Repr: "3"}},
						val:   Token{Class: Operator, Repr: "%"},
					},
					right: &Node{
						left:  &Node{val: Token{Class: Var, Repr: "i"}},
						right: &Node{val: Token{Class: Num, Repr: "5"}},
						val:   Token{Class: Operator, Repr: "%"},
					},
					val: Token{Class: Boolop, Repr: "&"},
				},
			},
		},
		{
			input: "print \"fizz\"",
			expected: []*Node{
				&Node{
					right: &Node{val: Token{Class: Str, Repr: "fizz"}},
					val:   Token{Class: Builtin, Repr: "print"},
				},
			},
		},
		{
			input: "if i % 3 == 0 print \"fizz\"",
			expected: []*Node{
				&Node{
					left: &Node{
						left: &Node{
							left:  &Node{val: Token{Class: Var, Repr: "i"}},
							right: &Node{val: Token{Class: Num, Repr: "3"}},
							val:   Token{Class: Operator, Repr: "%"},
						},
						right: &Node{val: Token{Class: Num, Repr: "0"}},
						val:   Token{Class: Boolop, Repr: "=="},
					},
					right: &Node{
						right: &Node{val: Token{Class: Str, Repr: "fizz"}},
						val:   Token{Class: Builtin, Repr: "print"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
			},
		},
		{
			input: "if i % 5 == 0 print \"buzz\"",
			expected: []*Node{
				&Node{
					left: &Node{
						left: &Node{
							left:  &Node{val: Token{Class: Var, Repr: "i"}},
							right: &Node{val: Token{Class: Num, Repr: "5"}},
							val:   Token{Class: Operator, Repr: "%"},
						},
						right: &Node{val: Token{Class: Num, Repr: "0"}},
						val:   Token{Class: Boolop, Repr: "=="},
					},
					right: &Node{
						right: &Node{val: Token{Class: Str, Repr: "buzz"}},
						val:   Token{Class: Builtin, Repr: "print"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
			},
		},
		{
			input: "if i < 100 goto 2",
			expected: []*Node{
				&Node{
					left: &Node{
						left: &Node{
							val: Token{Class: Var, Repr: "i"},
						},
						right: &Node{val: Token{Class: Num, Repr: "100"}},
						val:   Token{Class: Boolop, Repr: "<"},
					},
					right: &Node{
						right: &Node{val: Token{Class: Num, Repr: "2"}},
						val:   Token{Class: Builtin, Repr: "goto"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
			},
		},
		{
			input: `i = 0 
i = i + 1 
if i % 3 == 0 print "fizz" 
if i % 5 == 0 print "buzz" 
if i < 100 goto 2
`,
			expected: []*Node{
				&Node{
					left: &Node{
						val: Token{Class: Var, Repr: "i"},
					},
					right: &Node{
						val: Token{Class: Num, Repr: "0"},
					},
					val: Token{Class: Assignment, Repr: "="},
				},
				&Node{
					left: &Node{
						val: Token{Class: Var, Repr: "i"},
					},
					right: &Node{
						left:  &Node{val: Token{Class: Var, Repr: "i"}},
						right: &Node{val: Token{Class: Num, Repr: "1"}},
						val:   Token{Class: Operator, Repr: "+"},
					},
					val: Token{Class: Assignment, Repr: "="},
				},
				&Node{
					left: &Node{
						left: &Node{
							left:  &Node{val: Token{Class: Var, Repr: "i"}},
							right: &Node{val: Token{Class: Num, Repr: "3"}},
							val:   Token{Class: Operator, Repr: "%"},
						},
						right: &Node{val: Token{Class: Num, Repr: "0"}},
						val:   Token{Class: Boolop, Repr: "=="},
					},
					right: &Node{
						right: &Node{val: Token{Class: Str, Repr: "fizz"}},
						val:   Token{Class: Builtin, Repr: "print"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
				&Node{
					left: &Node{
						left: &Node{
							left:  &Node{val: Token{Class: Var, Repr: "i"}},
							right: &Node{val: Token{Class: Num, Repr: "5"}},
							val:   Token{Class: Operator, Repr: "%"},
						},
						right: &Node{val: Token{Class: Num, Repr: "0"}},
						val:   Token{Class: Boolop, Repr: "=="},
					},
					right: &Node{
						right: &Node{val: Token{Class: Str, Repr: "buzz"}},
						val:   Token{Class: Builtin, Repr: "print"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
				&Node{
					left: &Node{
						left: &Node{
							val: Token{Class: Var, Repr: "i"},
						},
						right: &Node{val: Token{Class: Num, Repr: "100"}},
						val:   Token{Class: Boolop, Repr: "<"},
					},
					right: &Node{
						right: &Node{val: Token{Class: Num, Repr: "2"}},
						val:   Token{Class: Builtin, Repr: "goto"},
					},
					val: Token{Class: Keyword, Repr: "if"},
				},
			},
		},
	}
	lexer := Lexer{}
	for _, test := range cases {
		t.Run(test.input, func(t *testing.T) {
			lexer.In = strings.NewReader(test.input)
			p := Parser{}
			tkns, _ := lexer.Lex()
			p.Tokens = tkns
			p.Parse()
			if !reflect.DeepEqual(p.Lines, test.expected) {
				t.Errorf("bad parse")
				for i := range test.expected {
					fmt.Println("expected")
					test.expected[i].PrintTree()
					fmt.Println("got")
					if i < len(p.Lines) {
						p.Lines[i].PrintTree()
					}
				}
			}
		})
	}
}
