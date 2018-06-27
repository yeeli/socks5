package main

import (
	"airship"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	file, err := ioutil.ReadFile("airship.yml")
	if err != nil {
		fmt.Errorf("err: %v", err)
	}
	app := &cli.App{
		Name: "airship",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start socks5 server",
				Action: func(c *cli.Context) error {
					conf := airship.ServerConfig{}
					yaml.Unmarshal(file, &conf)
					server := &airship.Server{
						Config: &conf,
					}
					server.Start()
					return nil
				},
			},
			{
				Name:  "client",
				Usage: "start socks5 client",
				Action: func(c *cli.Context) error {
					conf := airship.ClientConfig{}
					yaml.Unmarshal(file, &conf)
					client := &airship.Client{
						Config: &conf,
					}
					client.Start()
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
