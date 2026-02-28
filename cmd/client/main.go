package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	// default values
	DEFAULT_MESSAGE = "Hello World"
	BYTE_LIMIT      = 256
	LISTEN_PORT     = ":8080"
	PROTO           = "tcp"
	INTERACTIVE     = false
	// usage flag help (shown with `-h` flag)
	MESSAGE_USAGE     = "Pass a message (ASCII characters) to send to a server for encoding/decoding."
	LISTEN_PORT_USAGE = "Configure port for client to connect to."
	PROTO_USAGE       = "Configure protcool for client to connect with."
	BYTE_LIMIT_USAGE  = "Adjust bytes expected by/from the server."
	INTERACTIVE_USAGE = "Listen continously to user input."
)

type ClientRuntimeContext struct {
	message         string
	protocol        string
	listenPort      string
	byteLimit       int
	interactiveMode bool
}

func sendMessage(crContext ClientRuntimeContext, conn net.Conn) error {
	// push message to remote server

	// the write to the server here needs a `\n` to guarantee that the server
	// bothers to parse it (especially in one-shot that doesn't use newlines)
	// instead of just reading it infinitely
	_, err := fmt.Fprintf(conn, "%s\n", crContext.message)
	if err != nil {
		return fmt.Errorf("failed to write to remote server: %w", err)
	}

	// Prepare to accept the response of the server
	buffer := make([]byte, crContext.byteLimit)

	// read response from remote server
	numBytes, err := conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read server response: %w", err)
	}

	// server response
	fmt.Printf("Message returned: %s", string(buffer[:numBytes]))

	return nil
}

func runClient(crContext ClientRuntimeContext) error {
	// implement a TCP client as a Dialer?
	var dialer net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// open connection to remove server
	conn, err := dialer.DialContext(ctx, crContext.protocol, crContext.listenPort)
	if err != nil {
		return fmt.Errorf("failed to connect to remote server: %w", err)
	}

	if crContext.interactiveMode {
		// interactive mode:
		// take over the console, sending content of each line (up to `\n`) to server

		fmt.Println("You are now in interactive mode. Please type your message:")

		stdIn := bufio.NewScanner(os.Stdin)

		for stdIn.Scan() {
			// original `-message` contents don't matter, so we override it with the new input
			// and pass it to the sendMessage as usual
			crContext.message = stdIn.Text()
			err := sendMessage(crContext, conn)
			if err != nil {
				return err
			}
		}

		// check if scanner ran into any errors while scanning
		err := stdIn.Err()
		if err != nil {
			return fmt.Errorf("failed reading stdin: %w", err)
		}
	} else {
		// one-shot mode: read in `-message` and send that off to the server
		err := sendMessage(crContext, conn)
		if err != nil {
			return err
		}
	}

	// cleanup TCP connection
	err = conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

// parse command line arguments, with defaults
func parseArgs() *ClientRuntimeContext {
	message := flag.String("message", DEFAULT_MESSAGE, MESSAGE_USAGE)
	protcol := flag.String("proto", PROTO, PROTO_USAGE)
	listenPort := flag.String("port", LISTEN_PORT, LISTEN_PORT_USAGE)
	byteLimit := flag.Int("blimit", BYTE_LIMIT, BYTE_LIMIT_USAGE)
	interactiveMode := flag.Bool("interactive", INTERACTIVE, INTERACTIVE_USAGE)

	flag.Parse()

	return &ClientRuntimeContext{
		message:         *message,
		protocol:        *protcol,
		listenPort:      *listenPort,
		byteLimit:       *byteLimit,
		interactiveMode: *interactiveMode,
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
