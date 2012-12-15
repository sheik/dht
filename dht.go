/*
   Package dht implements a simple Distributed
   Hash Table.

   Check out the tests for usage examples.
*/
package dht

import (
	"fmt"
)

// k defines key length in bits
const k = 32

type Node struct {
	id     uint32
	next   *Node
	data   map[uint32]string
	finger map[uint32]*Node
}

// constructor for a new node
func NewNode(id uint32) *Node {
	return &Node{data: map[uint32]string{}}
}

// defines the distance between two keys
// using the kademlia xor metric
func distance(a, b uint32) uint32 {
	return a ^ b
}

// locate the node responsible for a key
func (start *Node) find(key uint32) *Node {
	current := start
	for distance(current.id, key) > distance(current.next.id, key) {
		current = current.next
	}
	return current
}

// Find the node responsible for the given
// key and return the value stored there
func (start *Node) lookup(key uint32) string {
	node := start.find(key)
	return node.data[key]
}

// Find the node responsible for the given
// key and store the value there
func (start *Node) store(key uint32, value string) {
	node := start.find(key)
	fmt.Println("Storing data in node with ID ", node.id)
	node.data[key] = value
}
