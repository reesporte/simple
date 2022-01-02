package simpl

// Stack is a basic stack
type Stack []*Node

// IsEmpty returns whether the stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds val to the stack
func (s *Stack) Push(val *Node) {
	*s = append(*s, val)
}

// Pop removes an element from the stack and returns it
func (s *Stack) Pop() *Node {
	if s.IsEmpty() {
		return nil
	}
	element := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return element
}

// Peek returns the top element of the stack but doesn't remove it
func (s *Stack) Peek() *Node {
	if s.IsEmpty() {
		return nil
	}
	return (*s)[len(*s)-1]
}
