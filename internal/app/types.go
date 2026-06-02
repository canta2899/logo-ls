package app

import (
	"github.com/canta2899/logo-ls/internal/icons"
	"github.com/canta2899/logo-ls/internal/inspect"
	"github.com/canta2899/logo-ls/pkg/fs"
)

// OpenDirIconString returns the colored open-directory glyph for a directory banner, or "" when icons are disabled.
func OpenDirIconString(showIcon bool) string {
	if !showIcon {
		return ""
	}
	d := icons.IconDef["diropen"]
	return d.GetColor() + d.GetGlyph() + "\033[0m" + " "
}

// FileEntry is a non-directory argument (e.g. logo-ls file.txt).
type FileEntry struct {
	Info    fs.FileInfo
	AbsPath string
}

func (f FileEntry) Name() string { return f.Info.Name() }

// DirectoryEntry holds an open directory handle so ReadDir is cheap.
type DirectoryEntry struct {
	File    fs.File
	AbsPath string
}

func (d *DirectoryEntry) Name() string { return d.File.Name() }
func (d *DirectoryEntry) Close() error { return d.File.Close() }

// Directory is the post-inspection, pre-render bundle for one directory.
type Directory struct {
	Info   *inspect.InspectedEntry
	Parent *inspect.InspectedEntry
	Files  []*inspect.InspectedEntry
	Dirs   []string
}
