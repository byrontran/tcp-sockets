package main

import (
	"fmt"
	"log"
	"net"
	"tcp-sockets/pkg/transform"
)

// is it fine that this is out of main scope?
const BYTE_LIMIT = 256

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
			fmt.Printf("Connection Accepted!\n")

			buffer := make([]byte, BYTE_LIMIT)
			num_bytes, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}

			if num_bytes > BYTE_LIMIT {
				fmt.Printf("This message exceeds the server's 256 character limit.\n")
				c.Close()
			}

			encoded_str := transform.Encode(string(buffer))
			fmt.Println(encoded_str)

			c.Close()
		}(conn)

	}
}
