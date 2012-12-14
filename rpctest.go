package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

// example rpc function that multiplies A and B
func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

// example rpc function that divides A by B
func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
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

// example rpc client
func test_client() {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatal("client: loadkeys: ", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	client, err := DialHTTPS("tcp", "localhost:1234", &config)

	if err != nil {
		log.Fatal("client: dial: ", err)
	}

	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
}

func main() {
	// start server startup code
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	// load TLS certificate
	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatal("server: loadkeys: ", err)
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	config.Rand = rand.Reader

	l, e := tls.Listen("tcp", ":1234", &config)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	fmt.Println("Server should be running...")
	// end server startup code

	test_client()
}
