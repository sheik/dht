package dht

import (
	"fmt"
	"math"
	"testing"
)

// tests the storing and lookup of keys
func TestLookup(t *testing.T) {
	first := SetupTestingCluster()
	for i := uint32(0); i < uint32(math.Pow(2, k)-10000010); i += 10000003 {
		first.store(i, fmt.Sprintf("TEST:%d", i))
	}

	for i := uint32(0); i < uint32(math.Pow(2, k)-10000010); i += 10000003 {
		val := first.lookup(i)
		expected := fmt.Sprintf("TEST:%d", i)
		if val != expected {
			t.Errorf("dht[%v] = %v, want %v", i, val, expected)
		}
	}
}

// sets up a test cluster to run tests on
func SetupTestingCluster() *Node {
	first := &Node{id: uint32(math.Pow(2, k) - 1), data: map[uint32]string{}}
	prev := first
	for i := uint32(math.Pow(2, k) - 1); i > 100000000; i -= 100000000 {
		new_prev := NewNode(i)
		new_prev.next = prev
		prev = new_prev
		fmt.Println("Creating node with ID: ", i)
	}
	first.next = prev
	return first
}
