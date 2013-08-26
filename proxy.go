// Adapted from https://gist.github.com/vmihailenco/1380352
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	localAddr = flag.String("local", "localhost:8888", "local address")
	remoteAddr = flag.String("remote", "localhost:9999", "remote address")
)

func proxyConn(conn *net.TCPConn) error {
	defer conn.Close()
	rAddr, err := net.ResolveTCPAddr("tcp", *remoteAddr)
	if err != nil {
		return err
	}

	rConn, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		return err
	}
	defer rConn.Close()

	errc := make(chan error)

	go copier(rConn, conn, "Local to Remote", errc)
	go copier(conn, rConn, "Remote to Local", errc)

	err1, err2 := <-errc, <-errc
	if err1 != nil {
		return err1
	}
	return err2
}

func handleConn(in <-chan *net.TCPConn) {
	for conn := range in {
		proxyConn(conn)
	}
}

func main() {
	flag.Parse()

	fmt.Printf("Local:  %v\n", *localAddr)
	fmt.Printf("Remote: %v\n\n", *remoteAddr)

	addr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	pending := make(chan *net.TCPConn)
	for i := 0; i < 5; i++ {
		go handleConn(pending)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("listener.AcceptTCP() -- %v\n", err)
			continue
		}
		pending <- conn
	}
}

func copier(w, r net.Conn, name string, errc chan error) {
	n, err := io.Copy(w, r)
	errc <- err
	fmt.Printf("%s: copied %d bytes\n", name, n)
}
