// Package inspect turns raw filesystem paths into InspectedEntry values that
// hold every piece of data the renderer needs. Inspection is the single
// place that touches the FS for per-file metadata; everything downstream
// reads from the resulting entries.
package inspect

import (
	iofs "io/fs"
	"time"

	"github.com/canta2899/logo-ls/icons"
)

// Kind classifies a filesystem entry.
type Kind int

const (
	KindFile Kind = iota
	KindDir
	KindSymlink
	KindPipe
	KindSocket
)

// InspectedEntry is the single source of truth for one filesystem entry.
// Populated once, in inspect.Inspect; never mutated downstream.
type InspectedEntry struct {
	Name    string // raw, as it appears on disk
	Ext     string // file extension (with leading dot) or "" for dotfiles/no-ext
	Base    string // Name with Ext stripped — kept so the renderer can compose Base+Ext+Indicator
	AbsPath string
	Kind    Kind
	Mode    iofs.FileMode // raw mode (not a formatted string)
	Size    int64
	ModTime time.Time

	// Long-listing extras. May be empty when the renderer doesn't need them.
	Inode     string
	HardLinks uint64
	Blocks    int64
	Owner     string // raw, unpadded
	Group     string // raw, unpadded
	HasXAttr  bool
	Sticky    bool // S_ISVTX is set
	StickyX   bool // sticky AND other-executable

	// Symlinks only. Populated when Kind == KindSymlink.
	LinkTarget   string // raw target as stored on disk (empty if EvalSymlinks not called)
	LinkResolved *InspectedEntry

	// Derived. Set by the inspector; safe to override in tests.
	Icon      *icons.IconInfo
	Indicator string
	GitStatus string
}

// IsDir reports whether the entry refers to a directory (post-symlink-follow
// semantics depend on how the entry was constructed).
func (e *InspectedEntry) IsDir() bool { return e.Kind == KindDir }

// IsSymlink reports whether the entry refers to a symlink.
func (e *InspectedEntry) IsSymlink() bool { return e.Kind == KindSymlink }

// kindFromMode picks the Kind that matches an os.FileMode.
func kindFromMode(m iofs.FileMode) Kind {
	switch {
	case m&iofs.ModeDir != 0:
		return KindDir
	case m&iofs.ModeSymlink != 0:
		return KindSymlink
	case m&iofs.ModeNamedPipe != 0:
		return KindPipe
	case m&iofs.ModeSocket != 0:
		return KindSocket
	default:
		return KindFile
	}
}
