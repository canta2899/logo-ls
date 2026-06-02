// Package format handles formatting and sorting logic
package format

import (
	"fmt"
	"os"
	"strings"

	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/model"
)

// compareName orders two names the way system ls does, honouring the active
// collation locale (LC_ALL -> LC_COLLATE -> LANG, read at call time). In the
// C/POSIX locale, comparison is by byte/ASCII order (uppercase before
// lowercase); in any other locale it falls back to case-insensitive order,
// matching the en_US.UTF-8 behaviour
func compareName(a, b string) bool {
	if cLocaleCollation() {
		return a < b
	}
	return strings.ToLower(a) < strings.ToLower(b)
}

// cLocaleCollation reports whether the effective collation locale is the
// C/POSIX locale (or unset), in which case byte-order comparison applies.
func cLocaleCollation() bool {
	loc := ""
	for _, k := range []string{"LC_ALL", "LC_COLLATE", "LANG"} {
		if v := os.Getenv(k); v != "" {
			loc = v
			break
		}
	}
	// Strip any codeset suffix, e.g. "C.UTF-8" → "C".
	if i := strings.IndexByte(loc, '.'); i >= 0 {
		loc = loc[:i]
	}
	return loc == "" || loc == "C" || loc == "POSIX"
}

// extGroupRank ranks an entry for -X (sort by extension) into one of three
// contiguous groups: 0 = extensionless non-dotfiles (dirs, LICENSE, Makefile),
// 1 = dotfiles (kept grouped, incl. "." and ".."), 2 = files with an extension.
func extGroupRank(name, ext string) int {
	if strings.HasPrefix(name, ".") {
		return 1
	}
	if ext == "" {
		return 0
	}
	return 2
}

func MainSort(a, b string) bool {
	aDot := strings.HasPrefix(a, ".")
	bDot := strings.HasPrefix(b, ".")

	// Dotfiles come before non-dotfiles
	if aDot != bDot {
		return aDot
	}

	// Within the same group, sort by name without the dot prefix (case-insensitive)
	if aDot {
		a = strings.TrimPrefix(a, ".")
	}
	if bDot {
		b = strings.TrimPrefix(b, ".")
	}
	return compareName(a, b)
}

// DotFileOrder checks whether dotfile grouping determines the order.
// Returns (result, decided): if decided is true, use result as the less value.
// "." and ".." are ordinary dotfiles here, sorted within the dotfile group by
// the caller's comparison rather than force-pinned to the top.
func DotFileOrder(a, b string) (bool, bool) {
	aDot := strings.HasPrefix(a, ".")
	bDot := strings.HasPrefix(b, ".")
	if aDot != bDot {
		return aDot, true
	}
	return false, false
}

// SetLessFunction is the custom less function to allow several sorting modes
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
			a := d.Files[i].Name + d.Files[i].Ext
			b := d.Files[j].Name + d.Files[j].Ext
			if res, ok := DotFileOrder(a, b); ok {
				return res
			}
			if d.Files[i].Size > d.Files[j].Size {
				return true
			} else if d.Files[i].Size == d.Files[j].Size {
				return MainSort(a, b)
			} else {
				return false
			}
		}
	case model.SortModTime:
		// sort by modification time, newest first
		d.LessFn = func(i, j int) bool {
			return d.Files[i].ModTime.After(d.Files[j].ModTime)
		}
	case model.SortExtension:
		// extensionless files/dirs first, then dotfiles (kept grouped), then
		// files sorted by extension. "." and ".." are ordinary dotfiles, sorted
		// within the dotfile group rather than force-pinned to the top.
		d.LessFn = func(i, j int) bool {
			a := d.Files[i].Name + d.Files[i].Ext
			b := d.Files[j].Name + d.Files[j].Ext
			ra := extGroupRank(d.Files[i].Name, d.Files[i].Ext)
			rb := extGroupRank(d.Files[j].Name, d.Files[j].Ext)
			if ra != rb {
				return ra < rb
			}
			if ra == 2 {
				if compareName(d.Files[i].Ext, d.Files[j].Ext) {
					return true
				}
				if compareName(d.Files[j].Ext, d.Files[i].Ext) {
					return false
				}
			}
			return MainSort(a, b)
		}
	case model.SortNatural:
		// natural sort of (version) numbers within text
		d.LessFn = func(i, j int) bool {
			a := d.Files[i].Name + d.Files[i].Ext
			b := d.Files[j].Name + d.Files[j].Ext
			if res, ok := DotFileOrder(a, b); ok {
				return res
			}
			return a < b
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

func GetOpenDirIcon() *icons.IconInfo {
	return icons.IconDef["diropen"]
}

func GetIcon(name, ext, indicator string) *icons.IconInfo {
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
	if indicator == "*" && i != nil {
		if i.GetGlyph() == "\uf723" {
			i = icons.IconDef["exe"]
		}
		i = i.AsExecutable()
	}

	if i == nil {
		i = icons.IconDef["file"]
	}

	return i
}

func GetFormattedSize(b int64, humanReadable bool) string {
	if !humanReadable {
		return fmt.Sprintf("%d", b)
	}

	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	size := float64(b) / float64(div)

	// Check if the size is a whole number
	if size == float64(int64(size)) {
		return fmt.Sprintf("%d%c", int64(size), "KMGTPE"[exp])
	}

	// Otherwise, keep one decimal place
	return fmt.Sprintf("%.1f%c", size, "KMGTPE"[exp])
}
