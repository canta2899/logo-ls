// Package ctw implements the column table writers used by the renderer.
package ctw

import (
	"bytes"
	"os"
	"strconv"

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
	if longMode {
		return NewLongCTW(10)
	}
	if oneFilePerLine {
		return NewLongCTW(4)
	}
	return NewStandardCTW(GetCustomTerminalWidth())
}

func GetCustomTerminalWidth() int {
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
