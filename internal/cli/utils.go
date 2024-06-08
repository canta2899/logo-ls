package cli

import (
	"io"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/canta2899/logo-ls/internal/app"
	"github.com/canta2899/logo-ls/internal/model"
	"github.com/mattn/go-colorable"
	"golang.org/x/term"
)

const standardTerminalWidth = 80

func GetCliApp() *app.App {
	command := GetConfig()
	writer := getCliWriter()
	logger := log.New(writer, "logo-ls: ", 0)
	terminalWidth := standardTerminalWidth

	if command.LongListingMode == model.LongListingNone {
		terminalWidth = getCustomTerminalWidth()
	}

	return &app.App{
		Config:        command,
		Writer:        writer,
		TerminalWidth: terminalWidth,
		Logger:        logger,
	}
}

func getCliWriter() io.Writer {
	var writer io.Writer

	writer = os.Stdout

	if runtime.GOOS == "windows" {
		writer = colorable.NewColorableStdout()
	}

	return writer
}

func getCustomTerminalWidth() int {
	// screen width for custom tw
	w, _, e := term.GetSize(int(os.Stdout.Fd()))

	if e != nil {
		return standardTerminalWidth
	}

	if w == 0 {
		// for systems that don’t support ‘TIOCGWINSZ’.
		w, _ = strconv.Atoi(os.Getenv("COLUMNS"))
	}

	return w
}
