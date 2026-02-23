package main

import (
	"context"
	"log"
	"net"
	"time"
)

func main() {
	// implement a TCP client as a Dialer?
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello World"))
	if err != nil {
		log.Fatal(err)
	}
}
