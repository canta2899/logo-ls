package ctw

import (
	"bytes"
)

type CTW interface {
	AddRow(args ...string)
	IconColor(c string)
	Flush(buf *bytes.Buffer)
	DisplayColor(b bool)
	GetGitColor(gitStatus string) string
	widthsSum(w [][4]int, p int) int
}
