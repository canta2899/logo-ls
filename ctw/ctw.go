package ctw

import (
	"bytes"
	"os"
	"strconv"

	"github.com/canta2899/logo-ls/model"
	"golang.org/x/term"
)

type CTW interface {
	AddRow(color string, args ...string)
	Flush(buf *bytes.Buffer)
	GetGitColor(gitStatus string) string
	widthsSum(w [][4]int, p int) int
}

const standardTerminalWidth = 80

func NewCTW(longMode, oneFilePerLine, icon bool) CTW {
	var out CTW

	if longMode {
		out = NewLongCTW(10)
	} else if oneFilePerLine {
		out = NewLongCTW(4)
	} else {
		out = NewStandardCTW(GetCustomTerminalWidth())
	}

	if !icon {
		model.OpenDirIcon = ""
	}

	return out
}

func GetCustomTerminalWidth() int {

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
