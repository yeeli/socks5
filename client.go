package airship

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	//"os"
	"time"
)

const (
	socks5Version = uint8(5)
	NoAuth        = uint8(0)
	UserPassAuth  = uint8(2)
)

type ClientConfig struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Chost string
	Cport string
}

type Client struct {
	Config *ClientConfig
}

func (c *Client) serverConn(conn net.Conn) error {
	defer conn.Close()

	// 本地不需要密码, 直接发送[5, 0]
	conn.Write([]byte{5, 0})

	buf := bufio.NewReader(conn)
	ver, _ := buf.Peek(3)

	sevHost := fmt.Sprintf("%s:%s", c.Config.Host, c.Config.Port)

	clientConn, err := net.Dial("tcp", sevHost)

	if err != nil {
		fmt.Errorf("err: %v", err)
	}

	defer clientConn.Close()

	//设置通行Channcel
	errCh := make(chan error, 2)
	//go proxy("send", conn, clientConn, errCh)

	info, _ := buf.Peek(8)
	ipLen := int(info[7])
	info2, _ := buf.Peek(ipLen + 10)
	host := string(info2[8 : ipLen+8])
	portInfo := info2[ipLen+8 : ipLen+10]
	port := (int(portInfo[0]) << 8) | int(portInfo[1])
	url := fmt.Sprintf("%s:%d", host, port)

	if ver[0] == socks5Version {
		socks := make([]byte, 3)
		if _, err := io.ReadAtLeast(buf, socks, len(socks)); err != nil {
			fmt.Errorf("Send data error: %v", err)
		}

		buf2 := bytes.NewBuffer(nil)
		buf2.Write([]byte{5, 1, 2})
		buf2.Write([]byte{1})
		userLen := byte(len(c.Config.User))
		buf2.Write([]byte{userLen})
		buf2.Write([]byte(c.Config.User))
		passLen := byte(len(c.Config.Pass))
		buf2.Write([]byte{passLen})
		buf2.Write([]byte(c.Config.Pass))
		data := io.MultiReader(buf2, buf)

		go proxy("send", url, clientConn, data, errCh)
	} else {
		go proxy("send", url, clientConn, buf, errCh)
	}

	clientBuf := bufio.NewReader(clientConn)
	socks2 := make([]byte, 4)
	if _, err := io.ReadAtLeast(clientBuf, socks2, len(socks2)); err != nil {
		fmt.Errorf("Recive data error: %v", err)
	}
	buf3 := bytes.NewBuffer(nil)
	data2 := io.MultiReader(buf3, clientBuf)

	go proxy("recive", url, conn, data2, errCh)

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

func proxy(data string, url string, dst io.Writer, src io.Reader, errCh chan error) {
	startAt := time.Now()
	//w := io.MultiWriter(os.Stdout, dst)
	len, err := io.Copy(dst, src)
	endAt := time.Now()
	subT := endAt.Sub(startAt)
	fmt.Printf("%s %s: %.2f KB complete in %v \n", url, data, float64(len)/float64(1024), subT)
	if tcpConn, ok := dst.(closeWriter); ok {
		tcpConn.CloseWrite()
	}
	errCh <- err
}

func (c *Client) Start() error {

	str := fmt.Sprintf("%s:%s", c.Config.Chost, c.Config.Cport)
	fmt.Println("start server", str)

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
		go c.serverConn(conn)
	}
	return nil
}
