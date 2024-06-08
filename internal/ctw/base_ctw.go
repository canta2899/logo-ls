package ctw

import (
	"strings"
)

type baseCtw struct {
	noColor string
	green   string
	brown   string
	empty   string
}

func newBaseCtw() *baseCtw {
	return &baseCtw{
		noColor: "\033[0m",
		green:   "\033[38;2;055;183;021m",
		brown:   "\033[38;2;192;154;107m",
		empty:   "\u0020",
	}
}

func (c *baseCtw) DisplayColor(b bool) {
	if !b {
		c.noColor = ""
		c.green = ""
		c.brown = ""
	}
}

func (c *baseCtw) GetGitColor(gitStatus string) string {
	switch strings.Trim(gitStatus, " ") {
	case "":
		return c.noColor
	case "U":
		return c.green
	default:
		return c.brown
	}
}

func (c *baseCtw) widthsSum(w [][4]int, p int) int {
	s := 0
	for _, v := range w {
		s += v[0] + v[1] + v[2] + v[3] + p
	}
	s -= p
	return s
}
