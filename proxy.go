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

	l2r := make(chan error)
	r2l := make(chan error)

	go copier(rConn, conn, "Local to Remote", l2r)
	go copier(conn, rConn, "Remote to Local", r2l)

	// Once one direction of copying fails, close both connections and
	// return

	select {
	case err = <-r2l:
		// fmt.Printf("r2l == %v\n", err)
		go func() {
			err = <-l2r
			// fmt.Printf("Latent l2r: %v\n", err)
		}()
	case err = <-l2r:
		// fmt.Printf("l2r == %v\n", err)
		go func() {
			err = <-r2l
			// fmt.Printf("Latent r2l: %v\n", err)
		}()
	}

	return err
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
	go func() {
		for {
			go func(conn *net.TCPConn) {
				log.Printf("New connection: %s <--> %s\n", conn.LocalAddr(),
					conn.RemoteAddr())
				if err := proxyConn(conn); err != nil {
					fmt.Printf("proxyConn: %v\n", err)
				}
			}(<-pending)
		}
	}()

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
