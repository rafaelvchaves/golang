package ch03

import (
	"io"
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	// If successful, OS exclusively assigns port 0 on IP 127.0.0.1 to
	// the listener.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	// Important to close the listener, otherwise memory leaks or deadlocks
	// can occur, since calls to Accept will block indefinitely.
	defer func() { _ = listener.Close() }()

	t.Logf("bound to %q", listener.Addr())

	// for {
	// 	// listener.Accept blocks until the listener detects an incoming connection
	// 	// and completes the TCP handshake.
	// 	conn, err := listener.Accept()

	// 	go func(c net.Conn) {
	// 		defer c.Close()
			

	// 	}

	// }

}

func TestDial(t *testing.T) {
	// Create a listener on a random port.
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	// Start listener in new goroutine
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			// listener.Accept blocks until the listener detects an incoming connection
			// and completes the TCP handshake.
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			// server handler function
			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					// Read data from client into buffer.
					n, err := c.Read(buf)
					if err != nil {
						// Read returns an io.EOF error upon receiving a FIN packet, so we
						// exit the handler.
						if err != io.EOF {
							t.Error(err)
						}
						return
					}

					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()
	
	// client side of connection.
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// initiate graceful termination from client side.
	conn.Close()
	<-done
	listener.Close()
	<-done
}