package ctw

import (
	"bytes"
	"fmt"
	"strings"
)

// LongCTW is the column writer for long mode listings.
type LongCTW struct {
	*baseCtw
	rows         [][]string
	columnWidths []int
	iconColors   []string
	iconPaddings []int
	// numCols = cols - 1; the icon and git status columns are handled separately.
	numCols int
}

func NewLongCTW(cols int) *LongCTW {
	return &LongCTW{
		baseCtw:      newBaseCtw(),
		numCols:      cols - 1,
		columnWidths: make([]int, cols),
		rows:         make([][]string, 0),
		iconColors:   make([]string, 0),
	}
}

func (l *LongCTW) AddRow(color string, iconPadding int, columns ...string) {
	if len(columns) != l.numCols+1 {
		return
	}

	for i, val := range columns {
		if len(val) > l.columnWidths[i] {
			l.columnWidths[i] = len(val)
		}
	}

	l.rows = append(l.rows, columns)
	l.iconColors = append(l.iconColors, color)
	l.iconPaddings = append(l.iconPaddings, iconPadding)
}

// Flush writes the table to buf, skipping zero-width columns.
func (l *LongCTW) Flush(buf *bytes.Buffer) {
	skipCols := make([]bool, len(l.columnWidths))
	for i, width := range l.columnWidths {
		if width == 0 {
			skipCols[i] = true
		}
	}

	l.columnWidths[l.numCols] = 1 // git column

	// Icon column: 1 for glyph + max padding across all rows.
	maxIconWidth := 1
	for _, p := range l.iconPaddings {
		if w := 1 + p; w > maxIconWidth {
			maxIconWidth = w
		}
	}
	l.columnWidths[l.numCols-2] = maxIconWidth

	for rowIdx, row := range l.rows {
		l.writeRow(buf, rowIdx, row, skipCols)
	}
}

func (l *LongCTW) writeRow(buf *bytes.Buffer, rowIdx int, row []string, skipCols []bool) {
	isFirstPrintedCol := true
	for colIdx, cellValue := range row {
		if skipCols[colIdx] {
			continue
		}
		if !isFirstPrintedCol {
			fmt.Fprint(buf, l.empty)
		}
		l.writeCell(buf, rowIdx, colIdx, cellValue, row, skipCols)
		isFirstPrintedCol = false
	}
	fmt.Fprintln(buf)
}

func (l *LongCTW) writeCell(buf *bytes.Buffer, rowIdx, colIdx int, cellValue string, row []string, skipCols []bool) {
	gitSkipped := skipCols[l.numCols]
	width := l.columnWidths[colIdx]

	switch {
	case colIdx == l.numCols-2:
		// Icon column: print the glyph and fill remaining width with spaces.
		fill := width - 1
		if fill < 0 {
			fill = 0
		}
		fmt.Fprintf(buf, "%s%1s%s%s", l.iconColors[rowIdx], cellValue, l.noColor, strings.Repeat(" ", fill))
	case colIdx >= l.numCols-1 && !gitSkipped:
		fmt.Fprintf(buf, "%s%-*s%s", l.GetGitColor(row[l.numCols]), width, cellValue, l.noColor)
	case gitSkipped && colIdx == l.numCols-1:
		fmt.Fprintf(buf, "%-*s", width, cellValue)
	case !gitSkipped && colIdx == 1:
		// permission column is left-aligned
		fmt.Fprintf(buf, "%-*s", width, cellValue)
	default:
		fmt.Fprintf(buf, "%*s", width, cellValue)
	}
}
