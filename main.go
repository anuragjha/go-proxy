package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

const version = "v0.0.0"

func main() {
	log.Print("Starting go-proxy:", version)

	l, err := net.Listen("tcp", "localhost:5001")
	if err != nil {
		log.Fatal("Listener: ", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Conn: ", err)
		}

		//go copyToStdErr(conn)
		go proxy(conn)
	}
}

func proxy(inConn net.Conn) {
	defer inConn.Close()

	outConn, err := net.Dial("tcp", "google.com:80")
	if err != nil {
		log.Printf("Dial conn finished in %v", err)
		return
	}
	defer outConn.Close()

	go io.Copy(outConn, inConn)
	io.Copy(inConn, outConn)
}

func copyToStdErr(conn net.Conn) {
	defer conn.Close()

	for {
		err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			log.Printf("SetReadDeadline finished with err: %v", err)
		}
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			log.Printf("Read from conn finished with err: %v", err)
			return
		}
		os.Stderr.Write(buf[:n])
	}

}
