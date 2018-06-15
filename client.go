package airship

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	socks5Version = uint8(5)
	NoAuth        = uint8(0)
	UserPassAuth  = uint8(2)
)

func serverConn(conn net.Conn) error {
	defer conn.Close()

	buf := bufio.NewReader(conn)

	ver, err := buf.Peek(1)

	clientConn, err := net.Dial("tcp", "127.0.0.1:6666")

	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	defer clientConn.Close()

	errCh := make(chan error, 2)
	if ver[0] == socks5Version {
		//buf.ReadByte()
		//buf.ReadByte()
		//buf.ReadByte()
		buf2 := bytes.NewBuffer(nil)
		buf2.Write([]byte{5, 1, 2})
		buf2.Write([]byte{1, 3, 'f', 'o', 'o', 3, 'b', 'a', 'r'})
		data := io.MultiReader(buf2, buf)
		go proxy(clientConn, data, errCh)
	} else {
		fmt.Println("1")
		go proxy(clientConn, buf, errCh)
	}
	go proxy(conn, clientConn, errCh)

	// Wait
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return e
		}
	}

	return nil
}

type closeWriter interface {
	CloseWrite() error
}

func proxy(dst io.Writer, src io.Reader, errCh chan error) {
	w := io.MultiWriter(os.Stdout, dst)
	len, err := io.Copy(w, src)
	fmt.Printf("\nproxy: %d \n", len)
	if tcpConn, ok := dst.(closeWriter); ok {
		tcpConn.CloseWrite()
	}
	errCh <- err
}

func Start() error {

	ln, err := net.Listen("tcp", "127.0.0.1:6677")

	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			return err
		}
		go serverConn(conn)
	}
	return nil
}
