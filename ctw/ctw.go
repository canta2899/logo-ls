package ctw

import (
	"bytes"

	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/model"
)

type CTW interface {
	AddRow(color string, args ...string)
	Flush(buf *bytes.Buffer)
	DisplayColor(b bool)
	GetGitColor(gitStatus string) string
	widthsSum(w [][4]int, p int) int
}

func NewCTW(longMode, oneFilePerLine, color, icon bool, terminalWidth int) CTW {
	var out CTW

	if longMode {
		out = NewLongCTW(10)
	} else if oneFilePerLine {
		out = NewLongCTW(4)
	} else {
		out = NewStandardCTW(terminalWidth)
	}

	if !color {
		out.DisplayColor(false)
		model.OpenDirIcon = icons.IconDef["diropen"].GetGlyph() + " "
	}

	if !icon {
		model.OpenDirIcon = ""
	}

	return out
}
