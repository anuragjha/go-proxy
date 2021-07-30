package main

import (
	"log"
	"net"
	"os"
	"time"
)

const version = "v0.0.0"

func main() {
	log.Print("Starting go-proxy:", version)

	l, err := net.Listen("tcp", "localhost:7878")
	if err != nil {
		log.Fatal("Listener: ", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Conn: ", err)
		}

		go copyToStdErr(conn)
	}
}

func copyToStdErr(conn net.Conn) {
	// Using Golang library to copy in Writer from Reader
	// n, err := io.Copy(os.Stderr, conn)
	//log.Printf("Copied %d bytes; finished with err %v", n, err)

	// Implemented custom version of copy where we set a conn idle timeout and also conn close
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
