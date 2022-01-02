package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/ueki-kazuki/golang-chaser/chaser"
)

const (
	ExitCodeOK int = iota
	ExitCodeError
)

type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	port := flag.Int("p", 2009, "CHaser server port. [2009]")
	host := flag.String("h", "127.0.0.1", "CHaser server hostname. [127.0.0.1]")
	name := flag.String("n", "User1", "Client name. [User1]")
	flag.Parse()

	client, err := chaser.NewClient(
		*name, *host, *port,
	)
	if err != nil {
		fmt.Println(err)
		return ExitCodeError
	}
	defer client.Close()

	for !client.GameSet {
		value, err := client.GetReady()
		if err != nil {
			break
		}
		if value[chaser.Down] != chaser.BLOCK {
			client.WalkDown()
		} else {
			client.WalkUp()
		}
	}
	return ExitCodeOK
}
