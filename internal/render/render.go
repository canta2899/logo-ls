// Package render owns the presentation layer: it consumes model entries
// and writes the rendered, column-aligned output. It is the only place
// that decides on column widths, indicator suffixes, time formatting and
// per-row ANSI codes.
package render

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/canta2899/logo-ls/ctw"
	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/model"
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

// EntryTimeFormatter renders an entry's modification time. Tests inject a
// fixed formatter so goldens don't drift.
type EntryTimeFormatter interface {
	Format(e *model.Entry) string
}

// Options carries everything the renderer needs that doesn't live on the
// individual entry.
type Options struct {
	Mode          Mode
	ShowIcon      bool
	ShowInode     bool
	ShowBlocks    bool
	HumanReadable bool
	TimeFormatter EntryTimeFormatter
}

// Render writes one directory's worth of entries to w in the selected mode.
func Render(w io.Writer, entries []*model.Entry, opts Options) {
	tw := ctw.NewCTW(opts.Mode == ModeLong, opts.Mode == ModeOneFilePerLine, opts.ShowIcon)
	for _, e := range entries {
		addRow(tw, e, opts)
	}
	buf := new(bytes.Buffer)
	tw.Flush(buf)
	_, _ = io.Copy(w, buf)
}

// addRow encodes the per-mode column layout. It's the only place that
// knows what the CTW's positional rows mean.
func addRow(tw ctw.CTW, e *model.Entry, opts Options) {
	switch opts.Mode {
	case ModeLong:
		tw.AddRow(
			e.Icon.GetColor(),
			blockSizeWithInode(e, opts),
			e.Mode,
			strconv.FormatUint(e.NumHardLinks, 10),
			e.Owner,
			e.Group,
			format.GetFormattedSize(e.Size, opts.HumanReadable),
			opts.TimeFormatter.Format(e),
			e.Icon.GetGlyph(),
			e.Name+e.Ext+e.Indicator,
			e.GitStatus,
		)
	case ModeOneFilePerLine:
		tw.AddRow(
			e.Icon.GetColor(),
			blockSizeWithInode(e, opts),
			e.Icon.GetGlyph(),
			e.Name+e.Ext+e.Indicator,
			e.GitStatus,
		)
	default:
		tw.AddRow(
			e.Icon.GetColor(),
			blockSizeWithInode(e, opts),
			e.Icon.GetGlyph(),
			e.Name+e.Ext+e.Indicator,
			e.GitStatus,
		)
	}
}

func blockSizeWithInode(e *model.Entry, opts Options) string {
	var parts []string
	if opts.ShowInode {
		parts = append(parts, e.InodeNumber)
	}
	if opts.ShowBlocks {
		parts = append(parts, format.GetFormattedSize(e.Blocks, opts.HumanReadable))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}
