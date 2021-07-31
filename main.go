package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var echoAddress = "localhost:5001"
var proxyAddress = "localhost:5000"

func main() {
	log.Printf("Starting echo: %s and proxy: %s", echoAddress, proxyAddress)

	// echoAddress Listen and Serve
	go listenAndServe("tcp", echoAddress, echo)
	// proxyAddress listen and serve to above listener
	go listenAndServe("tcp", proxyAddress, proxy)

	// Shutdown channel
	sigChan := make(chan os.Signal, 1)
	// Listen for close sig
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	exit := <-sigChan
	log.Printf("Shutting down with code %s", exit)
}

func listenAndServe(network string, address string, fn func(conn net.Conn)) {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("Listener: ", err)
	}
	log.Printf("%s / %s", network, address)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Conn: ", err)
		}

		go fn(conn)
	}
}

func proxy(inConn net.Conn) {
	defer inConn.Close()

	outConn, err := net.Dial("tcp", echoAddress)
	if err != nil {
		log.Printf("Dial conn finished in %v", err)
		return
	}
	defer outConn.Close()

	go io.Copy(outConn, inConn)
	io.Copy(inConn, outConn)
}

func echo(conn net.Conn) {
	defer conn.Close()

	for {
		err := conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		if err != nil {
			log.Printf("SetReadDeadline finished with err: %v", err)
		}
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("Read from conn finished with err: %v", err)
			return
		}

		os.Stderr.Write([]byte("Request: "))
		os.Stderr.Write(buf[:n])

		conn.Write([]byte("Response: "))
		conn.Write(buf[:n])
	}
}
