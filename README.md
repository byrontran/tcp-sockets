# CS576 - TCP Sockets (Programming Assignment 1)

Specs: Implement both a TCP client and TCP server.
- Any high level language with networking support (Initial choice: Go)


## Server:
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
