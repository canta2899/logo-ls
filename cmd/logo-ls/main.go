package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/canta2899/logo-ls/app"
	"github.com/canta2899/logo-ls/fs/osfs"
	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/internal/inspect/git"
	"github.com/canta2899/logo-ls/model"
	"github.com/mattn/go-colorable"
)

func main() {
	command := cli.GetConfigFromCli()

	var writer io.Writer = os.Stdout

	if runtime.GOOS == "windows" {
		writer = colorable.NewColorableStdout()
	}

	logger := log.New(writer, "logo-ls: ", 0)

	app := &app.App{
		Config:    command,
		Writer:    writer,
		Logger:    logger,
		ExitCode:  model.CodeOk,
		FS:        osfs.New(),
		GitReader: git.NewStatusReader(git.ExecPorcelain{}),
	}

	app.Run()

	app.Exit()
}
