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
					config := airship.ServerConfig{}
					yaml.Unmarshal(file, &config)
					airship.ServerStart(&config)
					return nil
				},
			},
			{
				Name:  "client",
				Usage: "start socks5 client",
				Action: func(c *cli.Context) error {
					config := airship.ClientConfig{}
					yaml.Unmarshal(file, &config)
					airship.Start(&config)
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
