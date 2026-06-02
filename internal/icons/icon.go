package icons

import "fmt"

// IconInfo is one entry in the icon set: a nerd-font glyph and an RGB color.
// IsExecutable swaps the color to a green hue for files marked executable.
type IconInfo struct {
	Glyph        string
	Color        [3]uint8 // RGB; (0,0,0) is black
	IsExecutable bool
}

// GetGlyph returns the icon glyph, or "" when the receiver is nil.
func (i *IconInfo) GetGlyph() string {
	if i == nil {
		return ""
	}
	return i.Glyph
}

// GetColor returns the ANSI 24-bit color escape for the icon. Executables
// always render in a fixed green regardless of the configured Color.
func (i *IconInfo) GetColor() string {
	if i == nil {
		return ""
	}
	if i.IsExecutable {
		return "\033[38;2;76;175;080m"
	}
	return fmt.Sprintf("\033[38;2;%03d;%03d;%03dm", i.Color[0], i.Color[1], i.Color[2])
}

// AsExecutable returns a copy of the icon with IsExecutable set, leaving the
// shared package-level icon definitions unchanged.
func (i *IconInfo) AsExecutable() *IconInfo {
	if i == nil {
		return nil
	}
	cp := *i
	cp.IsExecutable = true
	return &cp
}
