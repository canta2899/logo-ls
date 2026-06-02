package model

import (
	"github.com/canta2899/logo-ls/fs"
	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/internal/inspect"
)

// OpenDirIconString returns the colored "open directory" header used when
// printing a directory banner (e.g. `dir/:`). When showIcon is false the
// returned string is empty so the banner has no glyph.
func OpenDirIconString(showIcon bool) string {
	if !showIcon {
		return ""
	}
	d := icons.IconDef["diropen"]
	return d.GetColor() + d.GetGlyph() + "\033[0m" + " "
}

// FileEntry refers to an entity the CLI was asked to list directly (e.g.
// `logo-ls file.txt`).
type FileEntry struct {
	Info    fs.FileInfo
	AbsPath string
}

func (f FileEntry) Name() string { return f.Info.Name() }

// DirectoryEntry refers to a directory the CLI was asked to list. The File
// handle is kept open so ReadDir is cheap; Close releases it.
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
