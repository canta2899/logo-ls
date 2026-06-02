package ctw

import (
	"bytes"
	"fmt"
	"math"
)

// StandardCTW fits multiple columns into the terminal width.
// Each row has 4 sub-columns: size, icon, name+indicator, gitStatus.
type StandardCTW struct {
	*baseCtw
	rows          [][]string
	sizeWidths    []int
	nameWidths    []int
	gitWidths     []int
	iconColors    []string
	numCols       int
	showIcon      bool
	terminalWidth int
}

func NewStandardCTW(termW int) *StandardCTW {
	ctw := &StandardCTW{
		baseCtw:       newBaseCtw(),
		numCols:       3,
		terminalWidth: termW,
		rows:          make([][]string, 0),
		sizeWidths:    make([]int, 0),
		nameWidths:    make([]int, 0),
		gitWidths:     make([]int, 0),
		iconColors:    make([]string, 0),
	}
	return ctw
}

func (s *StandardCTW) AddRow(color string, args ...string) {
	if len(args) != s.numCols+1 {
		return
	}

	s.sizeWidths = append(s.sizeWidths, len(args[0]))
	s.nameWidths = append(s.nameWidths, len(args[2]))
	s.gitWidths = append(s.gitWidths, len(args[3]))

	if !s.showIcon && len(args[1]) > 0 {
		s.showIcon = true
	}

	s.rows = append(s.rows, args)
	s.iconColors = append(s.iconColors, color)
}

// Flush writes all rows in a multi-column layout that fits the terminal width.
func (s *StandardCTW) Flush(buf *bytes.Buffer) {
	rowCount := len(s.rows)
	if rowCount == 0 {
		return
	}

	columnPadding := 2
	optimalWidths := s.computeOptimalWidths(rowCount, columnPadding)
	columnsAcross := len(optimalWidths)
	rowsDown := int(math.Ceil(float64(rowCount) / float64(columnsAcross)))

	for rowIndex := range rowsDown {
		padding := columnPadding
		for colIndex := range columnsAcross {
			dataIndex := rowIndex + colIndex*rowsDown
			if dataIndex >= rowCount {
				continue
			}
			if colIndex == columnsAcross-1 {
				padding = 0
			}
			s.printRowCell(buf, dataIndex, optimalWidths[colIndex])
			fmt.Fprintf(buf, "%*s", padding, "")
		}
		fmt.Fprintln(buf)
	}
}

// computeOptimalWidths tries increasing column counts until the layout exceeds the terminal width,
// returning the widths [size, icon, name, gitStatus] for each column in the best layout found.
func (s *StandardCTW) computeOptimalWidths(rowCount, columnPadding int) [][4]int {
	intermediateWidths := make([][4]int, 0)
	var finalWidths [][4]int
	prevJumpValue := 0

	for {
		colCount := len(intermediateWidths) + 1
		intermediateWidths = append(intermediateWidths, [4]int{0, 0, 0, 0})

		jumpValue := int(math.Ceil(float64(rowCount) / float64(colCount)))
		if prevJumpValue == jumpValue {
			continue
		}
		prevJumpValue = jumpValue

		begin := 0
		end := jumpValue

		for i := 0; i < colCount && end <= rowCount; i++ {
			intermediateWidths[i] = s.calcSubColumnWidths(begin, end)
			begin, end = end, end+jumpValue
		}

		if end-jumpValue < rowCount {
			intermediateWidths[colCount-1] = s.calcSubColumnWidths(end-jumpValue, rowCount)
		}

		totalWidth := s.widthsSum(intermediateWidths, columnPadding)
		switch {
		case totalWidth > s.terminalWidth:
			if len(finalWidths) == 0 {
				finalWidths = make([][4]int, len(intermediateWidths))
				copy(finalWidths, intermediateWidths)
			}
			return finalWidths

		case totalWidth >= s.terminalWidth/2:
			finalWidths = make([][4]int, len(intermediateWidths))
			copy(finalWidths, intermediateWidths)
		}

		if colCount == rowCount {
			finalWidths = make([][4]int, len(intermediateWidths))
			copy(finalWidths, intermediateWidths)
			return finalWidths
		}
	}
}

func (s *StandardCTW) calcSubColumnWidths(begin, end int) [4]int {
	maxSizeWidth := 0
	maxNameWidth := 0
	maxGitWidth := 0

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
	if maxSizeWidth > 0 {
		result[0] = maxSizeWidth + 1
	}
	if s.showIcon {
		result[1] = 2
	}
	result[2] = maxNameWidth
	if maxGitWidth > 0 {
		result[3] = 2
	}
	return result
}

func (s *StandardCTW) printRowCell(buf *bytes.Buffer, rowIndex int, colSizes [4]int) {
	if colSizes[0] > 0 {
		fmt.Fprintf(buf, "%-*s%s", colSizes[0]-1, s.rows[rowIndex][0], s.empty)
	}

	if s.showIcon {
		fmt.Fprintf(buf, "%s%1s%s%s",
			s.iconColors[rowIndex],
			s.rows[rowIndex][1],
			s.noColor,
			s.empty,
		)
	}

	fmt.Fprintf(buf, "%s%-*s%s",
		s.GetGitColor(s.rows[rowIndex][3]),
		colSizes[2],
		s.rows[rowIndex][2],
		s.noColor,
	)

	if colSizes[3] > 0 {
		fmt.Fprintf(buf, "%s%s%1s%s",
			s.empty,
			s.GetGitColor(s.rows[rowIndex][3]),
			s.rows[rowIndex][3],
			s.noColor,
		)
	}
}
