package simpl

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

// Interpreter interprets simple ASTs
type Interpreter struct {
	Lines  *[]*Node
	Vars   map[string]interface{}
	w      io.Writer
	calls  Stack
	retval interface{}
}

// NewInterpreter creates a new Interpreter
func NewInterpreter(lines *[]*Node, writer io.Writer) Interpreter {
	i := Interpreter{Lines: lines, w: writer}
	i.Vars = make(map[string]interface{})
	return i
}

// Interpret interprets the ASTs `Lines` in the Interpreter
func (in *Interpreter) Interpret() interface{} {
	for i := len(*in.Lines) - 1; i >= 0; i-- {
		in.calls.Push((*in.Lines)[i])
	}
	for in.calls.Peek() != nil {
		cur := in.calls.Pop()
		//		cur.PrintTree()
		in.retval = in.eval(cur)
	}
	return in.retval
}

// eval evaluates a node
func (in *Interpreter) eval(n *Node) interface{} {
	if n == nil {
		return 0.0
	}
	switch n.val.Class {
	case Var:
		return in.Vars[n.val.Repr]
	case Str:
		return n.val.Repr
	case Num:
		v, _ := strconv.ParseFloat(n.val.Repr, 64)
		return v
	case Operator:
		left := in.eval(n.left)
		if left == nil {
			left = 0.0
		}
		right := in.eval(n.right)
		if right == nil {
			right = 0.0
		}
		switch n.val.Repr {
		case "+":
			switch left.(type) {
			case string:
				return stringAdd(left, right)
			case float64:
				switch right.(type) {
				case float64:
					return left.(float64) + right.(float64)
				case string:
					return stringAdd(left, right)
				}
			}
		case "-":
			return left.(float64) - right.(float64)
		case "*":
			return left.(float64) * right.(float64)
		case "/":
			return left.(float64) / right.(float64)
		case "%":
			return int(left.(float64)) % int(right.(float64))
		}
	case Boolop:
		left := toFloat64(in.eval(n.left))
		right := toFloat64(in.eval(n.right))
		switch n.val.Repr {
		// 0 is false
		case "&":
			return (left != 0) && (right != 0)
		case "|":
			return (left != 0) || (right != 0)
		case "==":
			return left == right
		case "!=":
			return left != right
		case ">":
			return left > right
		case "<":
			return left < right
		}
	case Builtin:
		right := in.eval(n.right)
		switch n.val.Repr {
		case "print":
			switch right := right.(type) {
			case string:
				fmt.Fprintf(in.w, ""+right+"")
			default:
				fmt.Fprintf(in.w, "%v", right)
			}
		case "goto":
			for i := len(*in.Lines) - 1; i > int(right.(float64))-2; i-- {
				in.calls.Push((*in.Lines)[i])
			}
		}
	case Assignment:
		variable := n.left.val.Repr
		right := in.eval(n.right)
		in.Vars[variable] = right
	case Keyword:
		if toFloat64(in.eval(n.left)) != 0 {
			in.eval(n.right)
		}
	default:
		log.Fatalf("cannot evaluate node of type %v", n.val.Class)
	}
	return nil
}

func toFloat64(val interface{}) float64 {
	switch val := val.(type) {
	case bool:
		if val {
			return 1
		}
		return 0
	case int:
		return float64(val)
	default:
		return val.(float64)
	}
}

func stringAdd(left, right interface{}) string {
	return fmt.Sprintf("%v", left) + fmt.Sprintf("%v", right)
}
