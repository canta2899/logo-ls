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
	disableColor()
}

const standardTerminalWidth = 80

func NewCTW(longMode, oneFilePerLine, noColor bool) CTW {
	var w CTW
	switch {
	case longMode:
		w = NewLongCTW(10)
	case oneFilePerLine:
		w = NewLongCTW(4)
	default:
		w = NewStandardCTW(GetCustomTerminalWidth())
	}
	if noColor {
		w.disableColor()
	}
	return w
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
