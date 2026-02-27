package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"tcp-sockets/pkg/transform"
)

// is it fine that this is out of main scope?
const (
	// default values
	BYTE_LIMIT  = 256
	LISTEN_PORT = ":8080"
	PROTO       = "tcp"
	// flag names
	ENCODE_FLAG = "encode"
	DECODE_FLAG = "decode"
	// usage flag help (shown with `-h` flag)
	TRANSFORM_USAGE   = "Configure server to encode or decode the passed message."
	LISTEN_PORT_USAGE = "Configure port for server to listen on."
	PROTO_USAGE       = "Configure protcool for server to listen to."
	BYTE_LIMIT_USAGE  = "Configure maximum bytes for server to accept."
)

type ServerRuntimeContext struct {
	transformMode string
	transformFunc func(string) string
	listenPort    string
	protocol      string
	byteLimit     int
}

// per-message client message handler
func handleUserResponse(ctx ServerRuntimeContext, c net.Conn, reportingChan chan error) {
	var err error
	defer func() {
		// if there is an error (i.e. function hit a != nil check and returned early), report it
		if err != nil {
			reportingChan <- err
		}

		// close the connection cleanly, or if that fails, report it
		closeErr := c.Close()
		if closeErr != nil {
			closeErr = fmt.Errorf("failed to close %s connection cleanly: %w", ctx.protocol, closeErr)
			reportingChan <- closeErr
			return
		}
	}()

	fmt.Printf("received connection from client: ")

	// read client message into buffer for parsing
	buffer := make([]byte, ctx.byteLimit+1)
	numBytes, readErr := c.Read(buffer)
	if readErr != nil {
		err = fmt.Errorf("failed to read message from buffer: %w", readErr)
		return
	}

	// because we check byteLimit + 1 into our buffer, we can check if
	// the client sent 257 bytes instead of 256
	// the prompt enforces "no more than 256 characters"
	if numBytes > ctx.byteLimit {
		err = fmt.Errorf("message exceeded server's byte limit (got: %d, wanted: %d)", numBytes, ctx.byteLimit)
		// while this may throw an error, it isn't important enough to resend. the client is doing something bad after all
		c.Write([]byte(err.Error()))
		return
	}

	// the buffer may include bytes that weren't filled, so we slice based on what we actually have
	bufferedString := string(buffer[:numBytes])
	receivedString := ctx.transformFunc((bufferedString))

	// server logging to tell what the client sent and what we encoded it as
	fmt.Printf("message: %s, encoded: %s\n", bufferedString, receivedString)

	// Instead of printing, write back to connection buffer?
	// sends the message back to the client to confirm that we got it
	_, writeErr := c.Write([]byte(receivedString))
	if writeErr != nil {
		err = fmt.Errorf("failed to write to connection buffer: %w", writeErr)
		return
	}
}

// long-lasting function to listen for messages until executation is interrupted
func runServer(srContext ServerRuntimeContext) error {
	var err error
	listener, err := net.Listen(srContext.protocol, srContext.listenPort)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}

	fmt.Printf("%s server listening on %s with transform directive: [%s]\n", srContext.protocol, srContext.listenPort, srContext.transformMode)

	reportingChannel := make(chan error, 10) // arbitrary buffer size for formatting goroutine errors

	// teardown logic for the listener
	defer func() {
		// I assume the chan will get closed on program exit, along with the goroutines.
		closeErr := listener.Close() // always close listener before process exits.
		if closeErr != nil {
			err = fmt.Errorf("failed to close %s listener: %w", srContext.protocol, closeErr)
		}
	}()

	// consumer for reporting channel to stdout
	go func() {
		for reportedErr := range reportingChannel {
			log.Printf("encountered issue with client request: %s\n", reportedErr)
		}
	}()

	// Loop infinitely for pending connections
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			var opErr net.OpError
			if errors.Is(acceptErr, &opErr) && opErr.Temporary() {
				// don't want to fall over if the error is transient
				// should cover timeouts and what-not (ECONNRESET, ECONNABORTED)
				// technically shouldn't use Temporary() since it's deprecated
				// but I don't know the alternative off the top of my head
				// https://cs.opensource.google/go/go/+/refs/tags/go1.26.0:src/net/net.go;l=552
				log.Printf("non-fatal issue occured during connection accept: %s", &opErr)
				continue
			} else {
				return fmt.Errorf("fatal issue occured while accepting request: %w", acceptErr)
			}
		}

		// in "net" package example, connection is handled in a
		// concurrent goroutine while the server continues listening
		// for more acceptions. Do we want this in our implementation?
		go handleUserResponse(srContext, conn, reportingChannel)
	}
}

// check if user passed a valid transform directive and return the associated function, else return error
func validateTransform(transformDirective string) (func(string) string, error) {
	switch transformDirective {
	case ENCODE_FLAG:
		return transform.Encode, nil
	case DECODE_FLAG:
		return transform.Decode, nil
	}
	return nil, fmt.Errorf("invalid transform provided: %s (expected: %s, %s)", transformDirective, ENCODE_FLAG, DECODE_FLAG)
}

// parse command line arguments, with defaults
func parseArgs() (*ServerRuntimeContext, error) {
	transformMode := flag.String("transform", ENCODE_FLAG, TRANSFORM_USAGE)
	listenPort := flag.String("port", LISTEN_PORT, LISTEN_PORT_USAGE)
	protocol := flag.String("proto", PROTO, PROTO_USAGE)
	byteLimit := flag.Int("blimit", BYTE_LIMIT, BYTE_LIMIT_USAGE)

	flag.Parse()

	// might as well cache the transform function that we will use, since we are here
	// saves us some cycles from needing an if-statement later
	transformFunc, err := validateTransform(*transformMode)
	if err != nil {
		return nil, err
	}

	return &ServerRuntimeContext{
		transformMode: *transformMode,
		transformFunc: transformFunc,
		listenPort:    *listenPort,
		protocol:      *protocol,
		byteLimit:     *byteLimit,
	}, nil
}

func main() {
	// get parameters for server executation
	context, err := parseArgs()
	if err != nil {
		err = fmt.Errorf("failed to parse flags: %w", err)
		log.Fatal(err)
	}

	// spin up the listener for the server
	err = runServer(*context)
	if err != nil {
		err = fmt.Errorf("%s connection closed due to error: %w", context.protocol, err)
		log.Fatal(err)
	}
}
