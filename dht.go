package main

import (
	"fmt"
	"math"
)

// k defines key length in bits
const k = 32

type Node struct {
	id     uint32
	next   *Node
	data   map[uint32]string
	finger map[uint32]*Node
}

func distance(a, b uint32) uint32 {
	switch {
	case a == b:
		return 0
	case a < b:
		return b - a
	}

	return uint32(math.Pow(2, k)) + (b - a)
}

// locate the node that is responsible for key
func (start *Node) find(key uint32) *Node {
	current := start
	for distance(current.id, key) > distance(current.next.id, key) {
		current = current.next
	}
	return current
}

// Find the node that holds the data and
// return it
func (start *Node) lookup(key uint32) string {
	node := start.find(key)
	return node.data[key]
}

// Find the node that holds the data and
// store it
func (start *Node) store(key uint32, value string) {
	node := start.find(key)
	fmt.Println("Storing data in node with ID ", node.id)
	node.data[key] = value
}

func setup_test_cluster() *Node {
	first := &Node{id: uint32(math.Pow(2, k) - 1), data: map[uint32]string{}}
	prev := first
	for i := uint32(math.Pow(2, k) - 1); i > 100000000; i -= 100000000 {
		new_prev := &Node{id: i, next: prev, data: map[uint32]string{}}
		prev = new_prev
		fmt.Println("Creating node with ID: ", i)
	}
	first.next = prev
	return first
}

func main() {
	// set up a cluster of nodes
	first := setup_test_cluster()

	for i := uint32(0); i < uint32(math.Pow(2, k)-10000010); i += 10000003 {
		first.store(i, fmt.Sprintf("Some test data -- %d", i))
	}

	for i := uint32(0); i < uint32(math.Pow(2, k)-10000010); i += 10000003 {
		fmt.Println("Key ", i, " holds: ", first.lookup(i))
	}
}
