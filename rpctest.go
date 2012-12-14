package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
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

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func test_client() {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatal("client: loadkeys: ", err)
	}
	fmt.Println("loaded client certs...")
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "localhost:1234", &config)
	if err != nil {
		log.Fatal("client: dial: ", err)
	}
	fmt.Println("Dialed in...")
	defer conn.Close()

	io.WriteString(conn, "CONNECT "+rpc.DefaultRPCPath+" HTTP/1.0\n\n")
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err != nil {
		log.Fatal("client: response: ", err)
	}
	fmt.Println("response: ", resp)

	client := rpc.NewClient(conn)
	fmt.Println("Created client: ", client)
	args := &Args{7, 8}
	var reply int
	fmt.Println("Making call...")
	err = client.Call("Arith.Multiply", args, &reply)
	fmt.Println("Returned?")
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("Made call!")
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
