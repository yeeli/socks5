package airship

import (
	"airship/socks5"
	"fmt"
)

type ServerConfig struct {
	Host string
	Port string
	User string
	Pass string
}

type Server struct {
	Config *ServerConfig
}

func (s *Server) Start() error {
	creds := socks5.StaticCredentials{
		s.Config.User: s.Config.Pass,
	}
	cator := socks5.UserPassAuthenticator{Credentials: creds}
	conf := &socks5.Config{
		AuthMethods: []socks5.Authenticator{cator},
	}

	/*conf := &socks5.Config{}*/
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	str := fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port)
	if err := server.ListenAndServe("tcp", str); err != nil {
		panic(err)
	}
	return nil
}
