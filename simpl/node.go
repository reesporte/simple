package simpl

import "fmt"

// Node is a node of an abstract syntax tree
type Node struct {
	left  *Node
	right *Node
	val   Token
}

// height gets the height of a given tree with root node `n`
func (n *Node) height() int {
	lheight := 0
	rheight := 0
	if n.left != nil {
		lheight = n.left.height()
	}
	if n.right != nil {
		rheight = n.right.height()
	}
	if lheight > rheight {
		return lheight + 1
	}
	return rheight + 1
}

// PrintTree prints the tree in breadth-first order
func (n *Node) PrintTree() {
	h := n.height()
	for i := 1; i < h+1; i++ {
		n.printLevel(i, "")
		fmt.Println()
	}
}

// printLevel prints the data at a given level of a tree
func (n *Node) printLevel(level int, side string) {
	if level == 1 {
		fmt.Printf("%v %v ", side, n.val)
	} else if level > 1 {
		if n.left != nil {
			n.left.printLevel(level-1, side+":L")
		}
		if n.right != nil {
			n.right.printLevel(level-1, side+":R")
		}
	}
}
