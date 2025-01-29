package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/canta2899/logo-ls/app"
	"github.com/canta2899/logo-ls/model"
	"github.com/mattn/go-colorable"
)

func main() {
	command := app.GetConfigFromCli()

	var writer io.Writer = os.Stdout

	if runtime.GOOS == "windows" {
		writer = colorable.NewColorableStdout()
	}

	logger := log.New(writer, "logo-ls: ", 0)

	app := &app.App{
		Config:   command,
		Writer:   writer,
		Logger:   logger,
		ExitCode: model.CodeOk,
	}

	app.Run()

	app.Exit()
}
