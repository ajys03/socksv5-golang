package main

import "fmt"
import "net"

const (
	BufferSize   = 1024
	SOCKSVersion = 0x05
)

func connectReq(conn net.Conn, addr byte) {

}

// currently no authentication
func handleMethodNegotiation(conn net.Conn) int {
	buf := make([]byte, BufferSize)

	// read SOCKS version and number of methods
	_, err := conn.Read(buf[:2])
	if err != nil || buf[0] != SOCKSVersion {
		fmt.Println("Invalid SOCKS version or number of methods")
		// clientConn.Close() -> return to handle method
		return -1
	}
	nMethods := int(buf[1])
	_, err = conn.Read(buf[:nMethods])
	if err != nil {
		fmt.Println("Error reading methods")
		// clientConn.Close()
		return -1
	}

	// Check for 'NO AUTHENTICATION REQUIRED' (X'00')
	for i := 0; i < nMethods; i++ {
		if buf[i] == 0x00 {
			// Send selected method (X'00' for 'NO AUTHENTICATION REQUIRED')
			_, err := conn.Write([]byte{0x05, 0x00})
			if err != nil {
				return 0
			}
			fmt.Println("Complete Method Handshake")
			return 0
		}
	}

	// If 'NO AUTHENTICATION REQUIRED' is not supported, send failure message
	_, err = conn.Write([]byte{0x05, 0xFF})
	fmt.Println("'NO AUTHENTICATION REQUIRED' is not supported")
	// clientConn.Close()
	return -1
}

func handleRequests(conn net.Conn) {
	buf := make([]byte, BufferSize)

	_, err := conn.Read(buf[:4])
	if err != nil || buf[0] != SOCKSVersion {
		fmt.Println("Invalid SOCKS version or request")
		// conn.Close()
		return
	}
	cmd := buf[1]
	addrType := buf[3]

	if cmd == 0x01 {
		fmt.Println("Received SOCKS CONNECT request.")
		connectReq(conn, addrType)
	}

	return
}

func handleConnection(conn net.Conn) {
	if handleMethodNegotiation(conn) == -1 {
		return
	}

	handleRequests(conn)

	fmt.Println("Finished Client Connection")
	return
}

func main() {
	fmt.Println("Starting Application...")

	ln, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("SOCKS 5 proxy server listening on port 1080...")
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ln)
	for {
		// ln.accept blocks till new connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println("Accepted connection from:", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}
