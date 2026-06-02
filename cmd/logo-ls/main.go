package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/canta2899/logo-ls/internal/app"
	"github.com/canta2899/logo-ls/pkg/fs/osfs"
	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/internal/icons"
	"github.com/canta2899/logo-ls/internal/inspect/git"
	"github.com/mattn/go-colorable"
)

func main() {
	command := cli.GetConfigFromCli()

	var writer io.Writer = os.Stdout

	if runtime.GOOS == "windows" {
		writer = colorable.NewColorableStdout()
	}

	logger := log.New(writer, "logo-ls: ", 0)

	var iconOverride *icons.Override
	switch {
	case command.NoIconOverride:
		// user opted out; leave nil
	case command.IconOverrideFile != "":
		ov, err := icons.LoadOverridesFromPath(command.IconOverrideFile)
		if err != nil {
			logger.Printf("ignoring icon overrides: %v\n", err)
		} else {
			iconOverride = ov
		}
	default:
		ov, err := icons.LoadOverrides()
		if err != nil {
			logger.Printf("ignoring icon overrides: %v\n", err)
		} else {
			iconOverride = ov
		}
	}

	app := &app.App{
		Config:       command,
		Writer:       writer,
		Logger:       logger,
		ExitCode:     cli.CodeOk,
		FS:           osfs.New(),
		GitReader:    git.NewStatusReader(git.ExecPorcelain{}),
		IconOverride: iconOverride,
	}

	app.Run()

	app.Exit()
}
