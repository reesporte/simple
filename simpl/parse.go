package simpl

// Parser holds the state needed for parsing
type Parser struct {
	Tokens      []Token
	Lines       []*Node
	operators   Stack
	operands    Stack
	assignments Stack
	keywords    Stack
}

// Parse parses the tokens into an AST
func (p *Parser) Parse() {
	for _, tkn := range p.Tokens {
		t := &Node{val: tkn}
		switch tkn.Class {
		case Str, Num, Var:
			p.operands.Push(t)
		case Assignment:
			p.handleToken(t, &p.assignments)
		case Boolop, Operator, Builtin:
			p.handleToken(t, &p.operators)
		case Keyword:
			p.handleToken(t, &p.keywords)
		case Paren:
			switch tkn.Repr {
			case "(":
				p.operators.Push(t)
			case ")":
				for p.operators.Peek().val.Repr != "(" {
					op := p.operators.Pop()
					if p.operands.Peek() != nil {
						op.right = p.operands.Pop()
					}
					if op.val.Class != Builtin { // builtins only have a right side
						if p.operands.Peek() != nil {
							op.left = p.operands.Pop()
						}
					}
					p.operands.Push(op)
				}
			}
		case Newline:
			p.emptyStacks()
		}
	}
	p.emptyStacks()
}

func (p *Parser) handleToken(t *Node, stack *Stack) {
	for stack.Peek() != nil && greaterPrecedence(stack.Peek(), t) {
		p.levelStack(stack)
	}
	stack.Push(t)
}

func (p *Parser) emptyStacks() {
	for p.operators.Peek() != nil {
		p.levelStack(&p.operators)
	}
	for p.assignments.Peek() != nil {
		p.levelStack(&p.assignments)
	}
	for p.keywords.Peek() != nil {
		p.levelStack(&p.keywords)
	}
	for p.operands.Peek() != nil {
		p.Lines = append(p.Lines, p.operands.Pop())
	}
}

// levelStack takes the thing from the operators stack, and gives it its args
func (p *Parser) levelStack(operators *Stack) {
	op := operators.Pop()
	if op.val.Repr == "(" {
		return
	}
	if p.operands.Peek() != nil {
		op.right = p.operands.Pop()
	}
	if op.val.Class != Builtin { // builtins only have a right side
		if p.operands.Peek() != nil {
			op.left = p.operands.Pop()
		}
	}
	p.operands.Push(op)
}

var precedences = map[string]int{
	"goto":  -1,
	"print": -1,
	">":     0,
	"<":     0,
	"==":    0,
	"!=":    0,
	"|":     1,
	"&":     1,
	"+":     2,
	"-":     2,
	"*":     3,
	"/":     3,
	"%":     3,
	"=":     4,
	"if":    5,
}

func greaterPrecedence(a, b *Node) bool {
	if a.val.Repr == "" {
		return false
	}
	return precedences[a.val.Repr] > precedences[b.val.Repr]
}
