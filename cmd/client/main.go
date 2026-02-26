package main

import (
  "fmt"
	"context"
	"flag"
	"log"
	"net"
	"time"
)

const byteLimit = 256

const defaultMessage = "Hello World"
const messageUsage = "Pass a message (ASCII characters, under 256 characters) to be messaged to the server for encoding/deoding."

var messageFlag = flag.String("message", defaultMessage, messageUsage)

func main() {
	flag.Parse()
	message := *messageFlag

	// implement a TCP client as a Dialer?
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}


  // Prepare to accept the response of the server
  buffer := make([]byte, byteLimit)

  _, err = conn.Read(buffer)
  if err != nil {
    log.Fatal(err)
  }
	defer conn.Close()

  fmt.Printf("Message returned: %s", string(buffer))
}
