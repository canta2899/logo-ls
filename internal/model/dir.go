// this file contain dir type definition
package model

import (
	"os"
	"time"

	"github.com/canta2899/logo-ls/assets"
)

// create the open dir icon
var OpenDirIcon = assets.Icon_Def["diropen"].GetColor(1) + assets.Icon_Def["diropen"].GetGlyph() + "\033[0m" + " "
var PathSeparator string = string(os.PathSeparator)

type FileInfo struct {
	os.FileInfo
	AbsPath string
}

type File struct {
	os.File
	AbsPath string
}

type Entry struct {
	Name, Ext, Indicator string
	ModTime              time.Time
	Size                 int64 // in bytes
	Mode                 string
	ModeBits             uint32
	Owner, Group         string // use syscall package
	Blocks               int64  // blocks required by the file multiply buy 512 to get block size
	// 'U'-> untracked file 'M'-> Modified file '●'-> modified dir ' '-> Not Updated/ not in git repo
	GitStatus string
	Icon      string
	IconColor string
}

type Dir struct {
	Info   *Entry
	Parent *Entry
	Files  []*Entry // all child files and dirs
	Dirs   []string // for recursion contain only child dirs
	LessFn func(int, int) bool
}

// define methods on *dir type only not on file type

// sorting functions
func (d *Dir) Len() int {
	return len(d.Files)
}

func (d *Dir) Swap(i, j int) {
	d.Files[i], d.Files[j] = d.Files[j], d.Files[i]
}

func (d *Dir) Less(i, j int) bool {
	return d.LessFn(i, j)
}
