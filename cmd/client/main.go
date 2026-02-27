package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	// default values
	DEFAULT_MESSAGE = "Hello World"
	BYTE_LIMIT      = 256
	LISTEN_PORT     = ":8080"
	PROTO           = "tcp"
	// usage flag help (shown with `-h` flag)
	MESSAGE_USAGE     = "Pass a message (ASCII characters) to be messaged to the server for encoding/deoding."
	LISTEN_PORT_USAGE = "Configure port for client to connect to."
	PROTO_USAGE       = "Configure protcool for client to connect with."
	BYTE_LIMIT_USAGE  = "Adjust bytes expected by/from the server."
)

type ClientRuntimeContext struct {
	message    string
	protocool  string
	listenPort string
	byteLimit  int
}

func runClient(crContext ClientRuntimeContext) error {
	// implement a TCP client as a Dialer?
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := dialer.DialContext(ctx, crContext.protocool, crContext.listenPort)
	if err != nil {
		return fmt.Errorf("failed to connect to remote server: %w", err)
	}

	_, err = conn.Write([]byte(crContext.message))
	if err != nil {
		return fmt.Errorf("failed to write to remote server: %w", err)
	}

	// Prepare to accept the response of the server
	buffer := make([]byte, crContext.byteLimit)

	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read server response: %w", err)
	}

	// server response
	fmt.Printf("Message returned: %s\n", string(buffer))

	// cleanup TCP connection
	err = conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

func parseArgs() *ClientRuntimeContext {
	message := flag.String("message", DEFAULT_MESSAGE, MESSAGE_USAGE)
	protcol := flag.String("proto", PROTO, PROTO_USAGE)
	listenPort := flag.String("port", LISTEN_PORT, LISTEN_PORT_USAGE)
	byteLimit := flag.Int("blimit", BYTE_LIMIT, BYTE_LIMIT_USAGE)

	flag.Parse()

	return &ClientRuntimeContext{
		message:    *message,
		protocool:  *protcol,
		listenPort: *listenPort,
		byteLimit:  *byteLimit,
	}
}

// oneshot function to send a message with command line to a remote server
func main() {
	ctx := parseArgs()
	err := runClient(*ctx)
	if err != nil {
		err := fmt.Errorf("client message send failed: %w", err)
		log.Fatal(err)
	}
}
