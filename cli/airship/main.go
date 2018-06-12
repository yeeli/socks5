package main

import (
	"airship"
	_ "fmt"
	_ "gopkg.in/urfave/cli.v2"
	"os"
)

func main() {
	args := os.Args
	if args[1] == "server" {
		airship.ServerStart()
	}

	if args[1] == "client" {
		airship.Start()
	}
}
