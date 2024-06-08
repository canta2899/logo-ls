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

func GetCliApp() *app.App {
	command := GetConfig()
	writer := getCliWriter()
	log.SetPrefix("logo-ls: ")
	log.SetFlags(0)

	terminalWidth := 80

	if command.LongListingMode == model.LongListingNone {
		terminalWidth = getCustomTerminalWidth()
	}

	return &app.App{
		Config:        command,
		Writer:        writer,
		TerminalWidth: terminalWidth,
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
		return 80
	}

	if w == 0 {
		// for systems that don’t support ‘TIOCGWINSZ’.
		w, _ = strconv.Atoi(os.Getenv("COLUMNS"))
	}

	return w
}
