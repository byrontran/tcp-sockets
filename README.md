# CS576 - TCP Sockets (Programming Assignment 1)

**TODOS:**

- Server should write back to the connection with the encoded/decoded message. Client should print.
- Put the client in an I/O loop for consistent message passing? Pass a flag for this mode perhaps.
- Do more work in the error handling (as opposed to logging Fatal every time).

### Notes for Teammates (Byron)

I am not familiar with Go best practices nor the language, so if you happen to be and
find style inconsistencies or poor practices anywhere feel free to correct or shout them
out!

Implementation is based on the examples I've found on the Go "net" package documentation.
The server follows the "Listener" example, while the client follows the "Dialer" example.

Link to the docs is here: https://pkg.go.dev/net

# Running Instructions

1. clone repo and cd into it
2. `go build ./...` from project root directory
3. in one terminal, run `go run ./cmd/server/. {--transform encode | decode}`
4. in another terminal, run `go run ./cmd/client/. {--message {PUT_MSG_HERE}}`
    - there are some issues with white space delimiting still, I'm not familiar at all
      with Go arg passing so there may be room for improvement here

> Can run the program without flags too. Default values are to encode the message, and for the
> client to send "Hello World" to the server.

# Program Specs

Implement both a TCP client and TCP server.
- Any high level language with networking support (Initial choice: Go)

## Server

- Server must accept a connection from a client, receives a text message with <= 256 chars.
- Server encodes message by replacing each char with next char in ASCII sequence.
  - (e.g. "Hello World" becomes "Ifmmp!Xpsme")
- Server responds with message converted to the encoded string.

## Client:

- Client must connect to server using same port.
- Client must pass the message to the server.
- Client receives encoded message from server and displays response.


### Optional Challenge

- Rather than encoding the message, allow server to decode when a flag is passed to it.
