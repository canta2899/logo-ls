package main

import (
	"github.com/canta2899/logo-ls/internal/cli"
)

func main() {
	app := cli.GetCliApp()
	app.Run()
	app.Exit()
}
