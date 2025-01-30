package ctw

import (
	"bytes"
	"fmt"
)

// LongCTW is a specialized column-table-writer for "long" listings.
// It contains logic to store data rows, compute column widths, and
// handle colorization for icons and Git status columns.
type LongCTW struct {
	*baseCtw

	// rows holds the table data. Each element in rows is itself a slice
	// of strings, one string per column.
	rows [][]string

	// columnWidths stores the computed width for each column.
	columnWidths []int

	// iconColors holds the color strings to be used for the icon column of each row.
	iconColors []string

	// numCols is the number of columns minus one (see usage).
	// Example: If the user says "5 columns," internally, numCols = 4.
	// This is a historical artifact—some columns are specially handled.
	numCols int
}

// NewLongCTW creates a new LongCTW that is set up to have `cols` columns.
//
// If cols = 5, for example, internally it will treat it as numCols = 4.
// This code historically assumes there's a special 'git column' plus an
// 'icon column'.
func NewLongCTW(cols int) *LongCTW {
	return &LongCTW{
		baseCtw:      newBaseCtw(),
		numCols:      cols - 1,
		columnWidths: make([]int, cols), // actual storage for column widths
		rows:         make([][]string, 0),
		iconColors:   make([]string, 0),
	}
}

// AddRow appends a new row to the table.
//
// The first parameter, iconColor, is the color that will be used to print
// the "icon column" for this row. The remaining arguments are the data
// for each column. The length of args must be exactly numCols+1.
func (l *LongCTW) AddRow(color string, columns ...string) {
	// Verify that the number of columns is correct
	if len(columns) != l.numCols+1 {
		return
	}

	// Update columnWidths if needed
	for i, val := range columns {
		if len(val) > l.columnWidths[i] {
			l.columnWidths[i] = len(val)
		}
	}

	l.rows = append(l.rows, columns)
	l.iconColors = append(l.iconColors, color)
}

// Flush writes the entire table to the provided buffer.
//
// It will skip columns that end up with zero-width, and it overrides
// certain columns (like the 'git column' and the 'icon column') to have
// explicit widths.
func (l *LongCTW) Flush(buf *bytes.Buffer) {
	// Determine which columns to skip (zero-width = skip)
	skipCols := make([]bool, len(l.columnWidths))
	for i, width := range l.columnWidths {
		if width == 0 {
			skipCols[i] = true
		}
	}

	// The last column (index l.numCols) is the "git column". Set width explicitly to 1.
	l.columnWidths[l.numCols] = 1

	// The "icon column" is l.numCols-2. Also set to 1 (or 2, if you need more space).
	l.columnWidths[l.numCols-2] = 1

	// Print each row
	for rowIdx, row := range l.rows {
		isFirstPrintedCol := true

		// Iterate columns in this row
		for colIdx, cellValue := range row {
			// Skip if the column has zero width
			if skipCols[colIdx] {
				continue
			}

			// If not the first column being printed, print separator
			if !isFirstPrintedCol {
				fmt.Fprint(buf, l.empty)
			}

			switch {
			// Icon column
			case colIdx == l.numCols-2:
				// Print with the icon color for this row
				fmt.Fprintf(buf, "%s%*s%s",
					l.iconColors[rowIdx],
					l.columnWidths[colIdx],
					cellValue,
					l.noColor,
				)

			// Git column (last column, if not skipped)
			case colIdx >= l.numCols-1 && !skipCols[l.numCols]:
				gitColor := l.GetGitColor(row[l.numCols])
				fmt.Fprintf(buf, "%s%-*s%s",
					gitColor,
					l.columnWidths[colIdx],
					cellValue,
					l.noColor,
				)

			// Git column is skipped, but we might still print the cellValue for some reason.
			// This block in the original code is to handle the case where the
			// last column is skipped, so it might still get printed differently.
			case skipCols[l.numCols] && colIdx == l.numCols-1:
				fmt.Fprintf(buf, "%-*s", l.columnWidths[colIdx], cellValue)

			// Normal column
			default:
				fmt.Fprintf(buf, "%*s", l.columnWidths[colIdx], cellValue)
			}

			isFirstPrintedCol = false
		}
		fmt.Fprintln(buf) // End of the row
	}
}
