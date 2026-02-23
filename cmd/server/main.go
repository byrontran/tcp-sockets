package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close() // always close listener before process exits.

	// Loop infinitely for pending connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// in "net" package example, connection is handled in a
		// concurrent goroutine while the server continues listening
		// for more acceptions. Do we want this in our implementation?
		go func(c net.Conn) {
			fmt.Printf("Connection Accepted!")
			c.Close()
		}(conn)

	}
}
