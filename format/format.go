package format

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/model"
)

func MainSort(a, b string) bool {
	switch a {
	case ".", "..":
	default:
		a = strings.TrimPrefix(a, ".")
	}
	switch b {
	case ".", "..":
	default:
		b = strings.TrimPrefix(b, ".")
	}
	return strings.ToLower(a) < strings.ToLower(b)
}

// Custom less functions
func SetLessFunction(d *model.Directory, sortMode model.SortMode) {
	switch sortMode {
	case model.SortAlphabetical:
		// sort by alphabetical order of name.ext
		d.LessFn = func(i, j int) bool {
			return MainSort(d.Files[i].Name+d.Files[i].Ext, d.Files[j].Name+d.Files[j].Ext)
		}
	case model.SortSize:
		// sort by file.Size, largest first
		d.LessFn = func(i, j int) bool {
			if d.Files[i].Size > d.Files[j].Size {
				return true
			} else if d.Files[i].Size == d.Files[j].Size {
				return MainSort(d.Files[i].Name+d.Files[i].Ext, d.Files[j].Name+d.Files[j].Ext)
			} else {
				return false
			}
		}
	case model.SortModTime:
		// sort by modification time, newest first
		// not sorting by alphabetical order because equality is quite rare
		d.LessFn = func(i, j int) bool {
			return d.Files[i].ModTime.After(d.Files[j].ModTime)
		}
	case model.SortExtension:
		// sort alphabetically by entry extension
		d.LessFn = func(i, j int) bool {
			if MainSort(d.Files[i].Ext, d.Files[j].Ext) {
				return true
			} else if strings.EqualFold(d.Files[i].Ext, d.Files[j].Ext) {
				return MainSort(d.Files[i].Name+d.Files[i].Ext, d.Files[j].Name+d.Files[j].Ext)
			} else {
				return false
			}
		}
	case model.SortNatural:
		// natural sort of (version) numbers within text
		d.LessFn = func(i, j int) bool {
			return d.Files[i].Name+d.Files[i].Ext < d.Files[j].Name+d.Files[j].Ext
		}
	case model.SortNone:
		fallthrough
	default:
		// no sorting
		d.LessFn = func(i, j int) bool {
			return i < j
		}
	}
}

// get indicator of the file
func GetIndicator(name string, isLongMode bool) (i string) {
	stats, err := os.Lstat(name)

	if err != nil {
		return ""
	}

	modebit := stats.Mode()

	switch {
	case modebit&os.ModeDir > 0:
		i = "/"
	case modebit&os.ModeNamedPipe > 0:
		i = "|"
	case modebit&os.ModeSymlink > 0:
		i = GetSymlinkIndicator(name, isLongMode)
	case modebit&os.ModeSocket > 0:
		i = "="
	case modebit&1000000 > 0:
		i = "*"
	}
	return i
}

func GetHardLinkCount(absPath string) uint64 {
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return 0
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}

	return uint64(stat.Nlink)
}

func GetInodeNumber(path string) string {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return ""
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)

	if !ok {
		return ""
	}

	return strconv.Itoa(int(stat.Ino))
}

func IsLink(name string) bool {
	stats, err := os.Lstat(name)

	if err != nil {
		return false
	}

	modebit := stats.Mode()

	return modebit&os.ModeSymlink > 0
}

func GetSymlinkIndicator(name string, isLongMode bool) string {
	if !isLongMode {
		return "@"
	}

	if s, err := filepath.EvalSymlinks(name); err == nil {
		return " ~> " + strings.Replace(s, os.Getenv("HOME"), "~", 1)
	}

	return ""
}

func GetIcon(name, ext, indicator string) (icon, color string) {
	var i *icons.IconInfo
	var ok bool

	switch indicator {
	case "/":
		i, ok = icons.IconDir[strings.ToLower(name+ext)]
		if ok {
			break
		}
		if len(name) == 0 || name[0] == '.' {
			i = icons.IconDef["hiddendir"]
			break
		}
		i = icons.IconDef["dir"]
	default:
		i, ok = icons.IconFileName[strings.ToLower(name+ext)]
		if ok {
			break
		}

		// a special admiration for goLang
		if ext == ".go" && strings.HasSuffix(name, "_test") {
			i = icons.IconSet["go-test"]
			break
		}

		t := strings.Split(name, ".")
		if len(t) > 1 && t[0] != "" {
			i, ok = icons.IconSubExt[strings.ToLower(t[len(t)-1]+ext)]
			if ok {
				break
			}
		}

		i, ok = icons.IconExt[strings.ToLower(strings.TrimPrefix(ext, "."))]
		if ok {
			break
		}

		if len(name) == 0 || name[0] == '.' {
			i = icons.IconDef["hiddenfile"]
			break
		}
		i = icons.IconDef["file"]
	}

	// change icon color if the file is executable
	if indicator == "*" {
		if i.GetGlyph() == "\uf723" {
			i = icons.IconDef["exe"]
		}
		i.MakeExe()
	}

	return i.GetGlyph(), i.GetColor(1)
}
