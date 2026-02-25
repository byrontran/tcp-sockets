package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"tcp-sockets/pkg/transform"
)

// is it fine that this is out of main scope?
const byteLimit = 256

const encodeFlag = "encode"
const decodeFlag = "decode"

const transformUsage = "Tell the TCP Server whether to encode or decode the passed message.\nEncode by default.\n"

var transformFlag = flag.String("transform", encodeFlag, transformUsage)

func main() {
	flag.Parse()
	transformMode := *transformFlag

	if transformMode != encodeFlag && transformMode != decodeFlag {
		fmt.Println("Server does not have valid instructions to handle messages from client.")
		fmt.Println("Please pass flags `--transform {encode | decode} to the server. If blank, server encodes by default.")

		os.Exit(1)
	}

	fmt.Printf("Starting the TCP server with mode: [%s]\n", transformMode)

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

			buffer := make([]byte, byteLimit)
			num_bytes, err := conn.Read(buffer)
			if err != nil {
				log.Fatal(err)
			}

			if num_bytes > byteLimit {
				fmt.Printf("This message exceeds the server's 256 character limit.\n")
				c.Close()
			}

			var resulting_str string

			if transformMode == encodeFlag {
				resulting_str = transform.Encode(string(buffer))
			} else {
				resulting_str = transform.Decode(string(buffer))
			}

			fmt.Println(resulting_str)

			c.Close()
		}(conn)

	}
}
