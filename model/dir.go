package model

import (
	"os"
	"sort"
	"time"

	"github.com/canta2899/logo-ls/icons"
)

var OpenDirIcon = icons.IconDef["diropen"].GetColor(1) + icons.IconDef["diropen"].GetGlyph() + "\033[0m" + " "
var PathSeparator string = string(os.PathSeparator)

type FileEntry struct {
	os.FileInfo
	AbsPath string
}

type DirectoryEntry struct {
	os.File
	AbsPath string
}

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
	Icon                 string
	IconColor            string
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
