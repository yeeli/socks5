package airship

import (
	"airship/socks5"
)

func ServerStart() {
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "127.0.0.1:6666"); err != nil {
		panic(err)
	}
}
