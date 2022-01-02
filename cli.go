package main

import (
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
	client, err := chaser.NewClient(
		"Test", "127.0.0.1", 2009,
	)
	if err != nil {
		fmt.Println(err)
		return ExitCodeError
	}
	defer client.Close()

	for {
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
