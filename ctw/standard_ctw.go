package ctw

import (
	"bytes"
	"fmt"
	"math"
)

// StandardCTW is a specialized column-table-writer for a "standard" listing.
// Unlike LongCTW, it tries to fit multiple columns of data in the available
// terminal width, along with optional icons and Git status columns.
//
// Each row has up to 4 sub-columns in the following order:
//  1. size
//  2. icon
//  3. name + extension + indicator
//  4. gitStatus
type StandardCTW struct {
	*baseCtw

	// Options:
	humanReadableSize bool
	showSize          bool

	// Data for each row. Each inner slice must have 4 columns:
	//   [0] -> size
	//   [1] -> icon
	//   [2] -> filename
	//   [3] -> git status
	rows [][]string

	// Tracking maximum widths for each sub-column, for each row:
	//   sizeWidths[i] is the length of rows[i][0]
	//   nameWidths[i] is the length of rows[i][2]
	//   gitWidths[i]  is the length of rows[i][3]
	//
	// The icon column has no "variable" width in the same sense, but
	// we track showIcon below to determine if we need to render the icon column.
	sizeWidths []int
	nameWidths []int
	gitWidths  []int

	// iconColors[i] is the ANSI color code for rows[i][1] (the "icon" column).
	iconColors []string

	// Number of columns (zero-based) within one row. Usually 3 means we are
	// expecting 4 sub-columns: size, icon, name, gitStatus.
	numCols int

	// Whether we have any icons in the dataset (i.e., if we need to print the icon column).
	showIcon bool

	// terminalWidth is the maximum width we can use for all columns combined.
	terminalWidth int
}

// NewStandardCTW creates a new StandardCTW, specifying the
// terminal width. The internal `numCols` is fixed at 3 (which
// yields 4 sub-columns: size, icon, name, gitStatus).
func NewStandardCTW(termW int) *StandardCTW {
	ctw := &StandardCTW{
		baseCtw:       newBaseCtw(),
		numCols:       3, // corresponds to 4 sub-columns
		terminalWidth: termW,
		rows:          make([][]string, 0),
		sizeWidths:    make([]int, 0),
		nameWidths:    make([]int, 0),
		gitWidths:     make([]int, 0),
		iconColors:    make([]string, 0),
	}
	return ctw
}

// AddRow appends a new row to the table, each row having exactly 4 pieces of data:
//
//	size, icon, name, gitStatus
//
// color is the ANSI color code for the icon column.
func (s *StandardCTW) AddRow(color string, args ...string) {
	// Verify correct number of columns
	if len(args) != s.numCols+1 { // s.numCols=3 => expects 4 pieces
		return
	}

	// Fill out max-width trackers for each sub-column
	s.sizeWidths = append(s.sizeWidths, len(args[0]))
	s.nameWidths = append(s.nameWidths, len(args[2]))
	s.gitWidths = append(s.gitWidths, len(args[3]))

	// If we haven't seen an icon yet, check if this row has one
	if !s.showIcon && len(args[1]) > 0 {
		s.showIcon = true
	}

	// Save the row data and the icon color
	s.rows = append(s.rows, args)
	s.iconColors = append(s.iconColors, color)
}

// Flush calculates how to best fit columns into the given terminal width,
// then prints them all to buf.
func (s *StandardCTW) Flush(buf *bytes.Buffer) {
	rowCount := len(s.rows)
	if rowCount == 0 {
		return
	}

	// We'll leave some spacing between columns:
	columnPadding := 2

	// Compute the best multi-column layout that fits in terminalWidth.
	optimalWidths := s.computeOptimalWidths(rowCount, columnPadding)

	// Once we have our final multi-column layout, figure out how many rows
	// we need for the entire data set. For example, if we ended up with
	// 3 columns across, we group the data in sets of rowCount/3, etc.
	columnsAcross := len(optimalWidths)
	rowsDown := int(math.Ceil(float64(rowCount) / float64(columnsAcross)))

	// Now, print all data in that multi-column layout
	for rowIndex := 0; rowIndex < rowsDown; rowIndex++ {
		// We'll reset spacing to columnPadding each time, except for the last column.
		padding := columnPadding

		for colIndex := 0; colIndex < columnsAcross; colIndex++ {
			// The actual index in s.rows we want to print is rowIndex + colIndex * rowsDown
			dataIndex := rowIndex + colIndex*rowsDown
			if dataIndex >= rowCount {
				continue // No more data in this "column"
			}

			// For the last column in a row, we do not add spacing after printing.
			if colIndex == columnsAcross-1 {
				padding = 0
			}

			s.printRowCell(buf, dataIndex, optimalWidths[colIndex])
			fmt.Fprintf(buf, "%*s", padding, "")
		}
		fmt.Fprintln(buf)
	}
}

// computeOptimalWidths iteratively tries to create "layouts" with
// increasing column counts, checking total required width each time,
// until it either exceeds the terminal width or it can place all
// entries in a single row.
//
// Each iteration returns a slice of `[4]int`, each `[4]int]` specifying
// the widths for size, icon, name, gitStatus.
func (s *StandardCTW) computeOptimalWidths(rowCount, columnPadding int) [][4]int {
	// This will be filled with sets of sub-column widths for each "layout"
	// attempt. For example, iw[0] is the sub-column widths for column #0,
	// iw[1] for column #1, etc.
	intermediateWidths := make([][4]int, 0)

	// We'll store the "best" set of widths that fit so far in `finalWidths`.
	var finalWidths [][4]int

	prevJumpValue := 0 // used to detect redundant calculations in loop

	for {
		colCount := len(intermediateWidths) + 1
		intermediateWidths = append(intermediateWidths, [4]int{0, 0, 0, 0})

		// jumpValue is how many rows go into each column, for this layout attempt
		jumpValue := int(math.Ceil(float64(rowCount) / float64(colCount)))
		if prevJumpValue == jumpValue {
			// If we haven't changed jumpValue from the previous iteration,
			// we're repeating ourselves, so skip to the next iteration.
			continue
		}
		prevJumpValue = jumpValue

		// b -> beginning index of the next chunk
		// e -> ending index of the next chunk
		begin := 0
		end := jumpValue

		// Fill each "column" chunk from begin..end, calculating the max sub-col widths
		for i := 0; i < colCount && end <= rowCount; i++ {
			intermediateWidths[i] = s.calcSubColumnWidths(begin, end)
			begin, end = end, end+jumpValue
		}

		// If the last column is not "complete" but still has leftover rows, compute them
		if end-jumpValue < rowCount {
			intermediateWidths[colCount-1] = s.calcSubColumnWidths(end-jumpValue, rowCount)
		}

		// Now compute the total width of the entire layout for colCount columns
		totalWidth := s.widthsSum(intermediateWidths, columnPadding)
		switch {
		case totalWidth > s.terminalWidth:
			// Once we exceed the terminal width, we stop. If we don't already have
			// a previously "good" layout, we keep the last best one we had.
			if len(finalWidths) == 0 {
				// We never found a layout that fit within the terminal, so we do a
				// fallback to the first one (like logo-ls -1).
				finalWidths = make([][4]int, len(intermediateWidths))
				copy(finalWidths, intermediateWidths)
			}
			return finalWidths

		case totalWidth >= s.terminalWidth/2:
			// If the layout's total width is more than half the terminal width,
			// we consider it a "good" layout and save it.
			finalWidths = make([][4]int, len(intermediateWidths))
			copy(finalWidths, intermediateWidths)
		}

		// If we have as many columns as rows, that means everything goes in one row.
		if colCount == rowCount {
			// We store the final widths to preserve them and then stop
			finalWidths = make([][4]int, len(intermediateWidths))
			copy(finalWidths, intermediateWidths)
			return finalWidths
		}
	}
}

// calcSubColumnWidths calculates the maximum needed widths (size, icon, name, gitstatus)
// for the rows in the slice [begin, end).
func (s *StandardCTW) calcSubColumnWidths(begin, end int) [4]int {
	maxSizeWidth := 0
	maxNameWidth := 0
	maxGitWidth := 0

	// For each row in [begin..end), update the max sub-column widths
	for i := begin; i < end; i++ {
		if s.sizeWidths[i] > maxSizeWidth {
			maxSizeWidth = s.sizeWidths[i]
		}
		if s.nameWidths[i] > maxNameWidth {
			maxNameWidth = s.nameWidths[i]
		}
		if s.gitWidths[i] > maxGitWidth {
			maxGitWidth = s.gitWidths[i]
		}
	}

	result := [4]int{0, 0, 0, 0}

	// If there's a size column to print
	if maxSizeWidth > 0 {
		// We add 1 to create a little space or alignment buffer
		result[0] = maxSizeWidth + 1
	}

	// If we are showing icons at all, we reserve 2 spaces for them
	if s.showIcon {
		result[1] = 2
	}

	// The filename column is assigned the maximum needed name width
	result[2] = maxNameWidth

	// For Git status, we also set it to 2 if there's anything to print
	if maxGitWidth > 0 {
		result[3] = 2
	}

	return result
}

// printRowCell writes a single row to the buffer based on the 4 sub-columns
// described by colSizes: [sizeColumnWidth, iconColumnWidth, nameColumnWidth, gitStatusColumnWidth].
func (s *StandardCTW) printRowCell(buf *bytes.Buffer, rowIndex int, colSizes [4]int) {
	// Sub-column 0: size
	if colSizes[0] > 0 {
		// %-*s left-justifies within colSizes[0]-1
		// then prints an extra space from s.empty
		fmt.Fprintf(buf, "%-*s%s", colSizes[0]-1, s.rows[rowIndex][0], s.empty)
	}

	// Sub-column 1: icon (only if showIcon is true)
	if s.showIcon {
		fmt.Fprintf(buf, "%s%1s%s%s",
			s.iconColors[rowIndex],
			s.rows[rowIndex][1],
			s.noColor,
			s.empty,
		)
	}

	// Sub-column 2: file name + indicator, color-coded by Git status
	fmt.Fprintf(buf, "%s%-*s%s",
		s.GetGitColor(s.rows[rowIndex][3]),
		colSizes[2],
		s.rows[rowIndex][2],
		s.noColor,
	)

	// Sub-column 3: Git status
	if colSizes[3] > 0 {
		fmt.Fprintf(buf, "%s%s%1s%s",
			s.empty,
			s.GetGitColor(s.rows[rowIndex][3]),
			s.rows[rowIndex][3],
			s.noColor,
		)
	}
}
