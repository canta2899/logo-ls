package model

import (
	"sort"
	"time"

	"github.com/canta2899/logo-ls/fs"
	"github.com/canta2899/logo-ls/icons"
)

var OpenDirIcon = icons.IconDef["diropen"].GetColor() + icons.IconDef["diropen"].GetGlyph() + "\033[0m" + " "

// PathSeparator is the path separator used by the active FS. The osfs
// implementation initializes this to the OS separator; tests override it.
var PathSeparator = "/"

type FileEntry struct {
	Info    fs.FileInfo
	AbsPath string
}

func (f FileEntry) Name() string { return f.Info.Name() }

type DirectoryEntry struct {
	File    fs.File
	AbsPath string
}

func (d *DirectoryEntry) Name() string { return d.File.Name() }
func (d *DirectoryEntry) Close() error { return d.File.Close() }

type Entry struct {
	Name, Ext, Indicator string
	ModTime              time.Time
	Size                 int64
	Mode                 string
	ModeBits             uint32
	NumHardLinks         uint64
	Owner, Group         string
	Blocks               int64
	GitStatus            string
	Icon                 *icons.IconInfo
	InodeNumber          string
}

type Directory struct {
	Info   *Entry
	Parent *Entry
	Files  []*Entry
	Dirs   []string
	LessFn func(int, int) bool
}

func (d *Directory) Len() int {
	return len(d.Files)
}

func (d *Directory) Swap(i, j int) {
	d.Files[i], d.Files[j] = d.Files[j], d.Files[i]
}

func (d *Directory) Less(i, j int) bool {
	return d.LessFn(i, j)
}

func (d *Directory) Sort(sortMode SortMode, reverse bool) {
	if sortMode != SortNone && reverse {
		sort.Sort(sort.Reverse(d))
	} else {
		sort.Sort(d)
	}
}
