package ctw

import (
	"bytes"
	"fmt"
)

// LongCTW is the column writer for long mode listings.
type LongCTW struct {
	*baseCtw
	rows         [][]string
	columnWidths []int
	iconColors   []string
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

func (l *LongCTW) AddRow(color string, columns ...string) {
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
}

// Flush writes the table to buf, skipping zero-width columns.
func (l *LongCTW) Flush(buf *bytes.Buffer) {
	skipCols := make([]bool, len(l.columnWidths))
	for i, width := range l.columnWidths {
		if width == 0 {
			skipCols[i] = true
		}
	}

	l.columnWidths[l.numCols] = 1   // git column
	l.columnWidths[l.numCols-2] = 1 // icon column

	for rowIdx, row := range l.rows {
		isFirstPrintedCol := true

		for colIdx, cellValue := range row {
			if skipCols[colIdx] {
				continue
			}

			if !isFirstPrintedCol {
				fmt.Fprint(buf, l.empty)
			}

			switch {
			case colIdx == l.numCols-2:
				fmt.Fprintf(buf, "%s%*s%s",
					l.iconColors[rowIdx],
					l.columnWidths[colIdx],
					cellValue,
					l.noColor,
				)

			case colIdx >= l.numCols-1 && !skipCols[l.numCols]:
				gitColor := l.GetGitColor(row[l.numCols])
				fmt.Fprintf(buf, "%s%-*s%s",
					gitColor,
					l.columnWidths[colIdx],
					cellValue,
					l.noColor,
				)

			case skipCols[l.numCols] && colIdx == l.numCols-1:
				fmt.Fprintf(buf, "%-*s", l.columnWidths[colIdx], cellValue)

			// permission column is left-aligned
			case !skipCols[l.numCols] && colIdx == 1:
				fmt.Fprintf(buf, "%-*s", l.columnWidths[colIdx], cellValue)

			default:
				fmt.Fprintf(buf, "%*s", l.columnWidths[colIdx], cellValue)
			}

			isFirstPrintedCol = false
		}
		fmt.Fprintln(buf)
	}
}
