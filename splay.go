/*
   Splay Tree with implicit index.

   https://neerc.ifmo.ru/wiki/index.php?title=Splay-%D0%B4%D0%B5%D1%80%D0%B5%D0%B2%D0%BE

   http://www.cs.cmu.edu/~sleator/papers/self-adjusting.pdf
*/

package main

import "fmt"

// Tree Node
type Node struct {
	value   int
	counter int
	left    *Node
	right   *Node
	parent  *Node
}

func ConstructNode(value int) Node {
	node := Node{
		value:   value,
		counter: 1,
		left:    nil,
		right:   nil,
		parent:  nil,
	}
	return node
}

func (node *Node) setLeft(leftNode *Node) {
	node.decrease(node.left)
	node.increase(leftNode)
	node.left = leftNode
	if leftNode != nil {
		leftNode.parent = node
	}
}

func (node *Node) setRight(rightNode *Node) {
	node.decrease(node.right)
	node.increase(rightNode)
	node.right = rightNode
	if rightNode != nil {
		rightNode.parent = node
	}
}

func (node *Node) increase(node2 *Node) {
	if node2 != nil {
		node.counter += node2.counter
	}
}

func (node *Node) decrease(node2 *Node) {
	if node2 != nil {
		node.counter -= node2.counter
	}
}

func (node *Node) isLeft() bool {
	return node.parent != nil && node.parent.left == node
}

func (node *Node) isRight() bool {
	return node.parent != nil && node.parent.right == node
}

func (node *Node) extract() {
	*node = ConstructNode(node.value)
}

type SplayTree struct {
	root  *Node
	step  int
	nodes []*Node
}

func ConstructSplayTree() SplayTree {
	st := SplayTree{
		root:  nil,
		step:  0,
		nodes: []*Node{},
	}
	return st
}

// Insert node with value ``value``
func (st *SplayTree) Insert(value int) int {
	node := ConstructNode(value)
	st.nodes = append(st.nodes, &node)
	pos := st.add(&node)
	return pos
}

// Remove node with position ``pos``
func (st *SplayTree) Remove(pos int) {
	node := st.FindByPos(pos)
	st.remove(node)
}

// Size of the tree
func (st *SplayTree) Size() int {
	if st.root != nil {
		return st.root.counter
	}
	return 0
}

// Set root of the tree
func (st *SplayTree) setRoot(node *Node) {
	st.root = node
	if st.root != nil {
		st.root.parent = nil
	}
}

// Rotate node
func (st *SplayTree) splay(node *Node) {
	if node == nil {
		return
	}
	for node.parent != nil {
		p := node.parent
		if p.parent == nil {
			// p is root
			st.stitch(node)
			break
		}

		if p.isLeft() && node.isLeft() ||
			p.isRight() && node.isRight() {
			st.straightStitch(node)
		} else {
			st.crossStitch(node)
		}
	}
	st.setRoot(node)
}

// Rotate node to root
func (st *SplayTree) stitch(node *Node) {
	p := node.parent
	if node.isLeft() {
		p.setLeft(node.right)
		node.setRight(p)
	} else {
		p.setRight(node.left)
		node.setLeft(p)
	}
	node.parent = nil
}

// Rotate left-left or right-right node.
func (st *SplayTree) straightStitch(node *Node) {
	p := node.parent
	g := p.parent

	if node.isLeft() {
		st.rotateRight(g)
		st.rotateRight(p)
	} else {
		st.rotateLeft(g)
		st.rotateLeft(p)
	}
}

// Rotate left-right node.
func (st *SplayTree) crossStitch(node *Node) {
	p := node.parent
	g := p.parent

	if node.isLeft() {
		st.rotateRight(p)
		st.rotateLeft(g)
	} else {
		st.rotateLeft(p)
		st.rotateRight(g)
	}
}

// Find the node and splay
func (st *SplayTree) find(node *Node) *Node {
	current := st.root
	for current != nil && current.value != node.value {
		if current.value < node.value {
			current = current.right
		} else {
			current = current.left
		}
	}
	st.splay(current)
	return current
}

// Find the node by value
func (st *SplayTree) FindByValue(value int) *Node {
	current := st.root
	for current != nil && current.value != value {
		if current.value < value {
			current = current.right
		} else {
			current = current.left
		}
	}
	return current
}

// Fint node by position and splay
func (st *SplayTree) FindByPos(pos int) *Node {
	node := st.root
	if node == nil || pos >= node.counter {
		fmt.Println("Position is out of range")
		return nil
	}
	for node != nil {
		if node.right != nil {
			if pos < node.right.counter {
				node = node.right
				continue
			}
			// skip all right nodes
			pos -= node.right.counter
		}
		if pos == 0 {
			break
		}
		// skip this node
		pos -= 1
		node = node.left
	}
	// self._splay(node)
	return node

}

// Merge trees.
func (st *SplayTree) merge(tree SplayTree) {
	if st.root == nil {
		st.setRoot(tree.root)
		return
	}
	node := st.MostRight()
	st.splay(node)
	node.setRight(tree.root)
}

// Split tree.
func (st *SplayTree) split(value int) (SplayTree, SplayTree) {
	node := st.root
	var next_node *Node

	left_tree := ConstructSplayTree()
	right_tree := ConstructSplayTree()

	if node == nil {
		return left_tree, right_tree
	}

	for {
		if node.value > value {
			next_node = node.left
		} else {
			next_node = node.right
		}
		if next_node == nil {
			break
		}
		node = next_node
	}

	st.splay(node)

	if node.value <= value {
		left_tree.setRoot(node)
		right_tree.setRoot(node.right)
		node.setRight(nil)
	} else {
		left_tree.setRoot(node.left)
		right_tree.setRoot(node)
		node.setLeft(nil)
	}

	st.setRoot(nil)
	return left_tree, right_tree
}

// Add node
func (st *SplayTree) add(node *Node) int {
	tree1, tree2 := st.split(node.value)
	st.setRoot(node)
	node.setLeft(tree1.root)
	node.setRight(tree2.root)
	pos := 0
	if node.right != nil {
		pos = node.right.counter
	}
	return pos
}

// Remove node
func (st *SplayTree) remove(node *Node) {
	st.splay(node)
	st.setRoot(node.left)
	right_tree := ConstructSplayTree()
	right_tree.setRoot(node.right)
	st.merge(right_tree)
	node.extract()
}

// Most left Leaf of the node
func (st *SplayTree) MostLeft() *Node {
	size := st.root.counter
	return st.FindByPos(size - 1)

}

// Most right Leaf of the node
func (st *SplayTree) MostRight() *Node {
	return st.FindByPos(0)
}

// Rotate edge from node to right child counter -clockwise.
func (st *SplayTree) rotateLeft(node *Node) {
	child := node.right

	counter_node := node.counter
	counter_child := child.counter

	p := node.parent
	if p != nil {
		counter_p := p.counter
		if node.isLeft() {
			p.setLeft(child)
		} else if node.isRight() {
			p.setRight(child)
		}
		p.counter = counter_p
	}

	child.parent = p

	// child.left, node.right = node, child.left
	grand_child := child.left
	child.setLeft(node)
	node.setRight(grand_child)

	node.counter = counter_node - counter_child

	if grand_child != nil {
		node.counter += grand_child.counter
		counter_child -= grand_child.counter
	}

	child.counter = counter_child + node.counter
}

// Rotate edge from node to left child clockwise.
func (st *SplayTree) rotateRight(node *Node) {
	child := node.left

	counter_node := node.counter
	counter_child := child.counter

	p := node.parent
	if p != nil {
		counter_p := p.counter
		if node.isLeft() {
			p.left = child
		} else if node.isRight() {
			p.right = child
		}
		p.counter = counter_p
	}

	child.parent = p

	// child.right, node.left = node, child.right
	grand_child := child.right
	child.setRight(node)
	node.setLeft(grand_child)

	node.counter = counter_node - counter_child

	if grand_child != nil {
		node.counter += grand_child.counter
		counter_child -= grand_child.counter
	}

	child.counter = counter_child + node.counter
}
