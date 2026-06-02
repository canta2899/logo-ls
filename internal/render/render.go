// Package render owns the presentation layer: it consumes
// *inspect.InspectedEntry values and writes the rendered, column-aligned
// output. It is the only place that decides on column widths, indicator
// suffixes, time formatting and per-row ANSI codes.
package render

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/canta2899/logo-ls/internal/render/columns"
	"github.com/canta2899/logo-ls/internal/inspect"
)

// Mode selects which renderer to use.
type Mode int

const (
	// ModeShort produces a multi-column layout sized to the terminal.
	ModeShort Mode = iota
	// ModeOneFilePerLine writes one entry per line, no extra columns.
	ModeOneFilePerLine
	// ModeLong renders the `ls -l`-style multi-column long listing.
	ModeLong
)

// TimeFormatter renders a single time.Time. Tests inject a deterministic
// formatter so goldens don't drift.
type TimeFormatter interface {
	Format(t *time.Time) string
}

// Options carries everything the renderer needs that doesn't live on the
// individual entry.
type Options struct {
	Mode          Mode
	ShowIcon      bool
	ShowInode     bool
	ShowBlocks    bool
	HumanReadable bool
	TimeFormatter TimeFormatter
}

// Render writes one directory's worth of entries to w in the selected mode.
func Render(w io.Writer, entries []*inspect.InspectedEntry, opts Options) {
	tw := ctw.NewCTW(opts.Mode == ModeLong, opts.Mode == ModeOneFilePerLine, opts.ShowIcon)
	for _, e := range entries {
		addRow(tw, e, opts)
	}
	buf := new(bytes.Buffer)
	tw.Flush(buf)
	_, _ = io.Copy(w, buf)
}

func addRow(tw ctw.CTW, e *inspect.InspectedEntry, opts Options) {
	displayName := e.Base + e.Ext + e.Indicator
	switch opts.Mode {
	case ModeLong:
		tw.AddRow(
			e.Icon.GetColor(),
			blockSizeWithInode(e, opts),
			inspect.ModeString(e.Mode, e.Sticky, e.StickyX, e.HasXAttr),
			strconv.FormatUint(hardLinks(e), 10),
			e.Owner,
			paddedGroup(e.Group),
			formatSize(e.Size, opts.HumanReadable),
			opts.TimeFormatter.Format(&e.ModTime),
			e.Icon.GetGlyph(),
			displayName,
			e.GitStatus,
		)
	case ModeOneFilePerLine:
		tw.AddRow(
			e.Icon.GetColor(),
			blockSizeWithInode(e, opts),
			e.Icon.GetGlyph(),
			displayName,
			e.GitStatus,
		)
	default:
		tw.AddRow(
			e.Icon.GetColor(),
			blockSizeWithInode(e, opts),
			e.Icon.GetGlyph(),
			displayName,
			e.GitStatus,
		)
	}
}

func hardLinks(e *inspect.InspectedEntry) uint64 {
	if e.HardLinks == 0 {
		return 1
	}
	return e.HardLinks
}

// paddedGroup wraps a group name in the legacy " %v  " padding the long
// listing's column widths were tuned for. Empty strings stay empty so the
// column collapses when -G is set.
func paddedGroup(g string) string {
	if g == "" {
		return ""
	}
	return fmt.Sprintf(" %v  ", g)
}

func blockSizeWithInode(e *inspect.InspectedEntry, opts Options) string {
	var parts []string
	if opts.ShowInode {
		parts = append(parts, e.Inode)
	}
	if opts.ShowBlocks {
		parts = append(parts, formatSize(e.Blocks, opts.HumanReadable))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}
