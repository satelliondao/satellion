package main

import (
	"github.com/satelliondao/satellion/cli"
)

func main() {
	cli.SetupCommands()
	cli.Execute()
}
