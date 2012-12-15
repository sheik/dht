/*
   Package dht implements a simple Distributed
   Hash Table.

   Check out the tests for usage examples.
*/
package dht

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/rpc"
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
	node.data[key] = value
}

// DialHTTPS connects to an HTTPS RPC server at the specified network address
// listening on the default HTTP RPC path using tls.Config config
func DialHTTPS(network, address string, config *tls.Config) (*rpc.Client, error) {
	return DialHTTPSPath(network, address, rpc.DefaultRPCPath, config)
}

// DialHTTPSPath connects to an HTTPS RPC server
// at the specified network address and path using tls.Config config
func DialHTTPSPath(network, address, path string, config *tls.Config) (*rpc.Client, error) {
	var err error
	conn, err := tls.Dial(network, address, config)
	if err != nil {
		return nil, err
	}
	io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	// Require successful HTTP response
	// before switch to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil && resp.StatusCode == http.StatusOK {
		return rpc.NewClient(conn), nil
	}
	if err == nil {
		err = errors.New("unexpected HTTP response: " + resp.Status)
	}
	conn.Close()
	return nil, &net.OpError{
		Op:   "dial-http",
		Net:  network + " " + address,
		Addr: nil,
		Err:  err,
	}
}
