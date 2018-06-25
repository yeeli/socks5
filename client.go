package airship

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

type ClientConfig struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Cport string
}

const (
	socks5Version = uint8(5)
	NoAuth        = uint8(0)
	UserPassAuth  = uint8(2)
)

func serverConn(c *ClientConfig, conn net.Conn) error {
	defer conn.Close()

	buf := bufio.NewReader(conn)
	ver, _ := buf.Peek(3)

	sevHost := fmt.Sprintf("%s:%s", c.Host, c.Port)

	clientConn, err := net.Dial("tcp", sevHost)

	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	defer clientConn.Close()

	// 本地不需要密码, 直接发送[5, 0]
	conn.Write([]byte{5, 0})

	errCh := make(chan error, 2)
	go proxy(conn, clientConn, errCh)

	if ver[0] == socks5Version {
		socks := make([]byte, 3)
		if _, err := io.ReadAtLeast(buf, socks, len(socks)); err != nil {
			fmt.Println(err)
		}
		buf2 := bytes.NewBuffer(nil)
		buf2.Write([]byte{5, 1, 2})
		buf2.Write([]byte{1})
		userLen := byte(len(c.User))
		buf2.Write([]byte{userLen})
		buf2.Write([]byte(c.User))
		passLen := byte(len(c.Pass))
		buf2.Write([]byte{passLen})
		buf2.Write([]byte(c.Pass))
		data := io.MultiReader(buf2, buf)
		go proxy(clientConn, data, errCh)
	} else {
		go proxy(clientConn, buf, errCh)
	}

	clientBuf := bufio.NewReader(clientConn)
	socks2 := make([]byte, 4)
	if _, err := io.ReadAtLeast(clientBuf, socks2, len(socks2)); err != nil {
		fmt.Println(err)
	}
	buf3 := bytes.NewBuffer(nil)
	data2 := io.MultiReader(buf3, clientBuf)

	go proxy(conn, data2, errCh)

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
	//w := io.MultiWriter(os.Stdout, dst)
	_, err := io.Copy(dst, src)
	//fmt.Printf("proxy: %d \n", len)
	if tcpConn, ok := dst.(closeWriter); ok {
		tcpConn.CloseWrite()
	}
	errCh <- err
}

func Start(c *ClientConfig) error {

	str := fmt.Sprintf("%s:%s", c.Host, c.Cport)

	ln, err := net.Listen("tcp", str)

	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			return err
		}
		go serverConn(c, conn)
	}
	return nil
}
