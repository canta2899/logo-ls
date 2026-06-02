// Package sort orders []*inspect.InspectedEntry by the user-selected sort
// mode. It is pure: no FS access, no globals. Locale handling reads the
// LC_* env at call time so the test harness's LC_ALL=C still works.
package sort

import (
	"os"
	stdsort "sort"
	"strings"

	"github.com/canta2899/logo-ls/internal/inspect"
	"github.com/canta2899/logo-ls/model"
)

// Sort sorts entries in place by mode. When reverse is true (and mode is
// not SortNone) the result is reversed.
func Sort(entries []*inspect.InspectedEntry, mode model.SortMode, reverse bool) {
	less := lessFn(entries, mode)
	if mode == model.SortNone {
		return
	}
	if reverse {
		stdsort.Slice(entries, func(i, j int) bool { return less(j, i) })
		return
	}
	stdsort.Slice(entries, less)
}

func lessFn(entries []*inspect.InspectedEntry, mode model.SortMode) func(i, j int) bool {
	switch mode {
	case model.SortAlphabetical:
		return func(i, j int) bool {
			return mainSort(entries[i].Name, entries[j].Name)
		}
	case model.SortSize:
		return func(i, j int) bool {
			a, b := entries[i].Name, entries[j].Name
			if res, ok := dotFileOrder(a, b); ok {
				return res
			}
			if entries[i].Size > entries[j].Size {
				return true
			}
			if entries[i].Size == entries[j].Size {
				return mainSort(a, b)
			}
			return false
		}
	case model.SortModTime:
		return func(i, j int) bool {
			return entries[i].ModTime.After(entries[j].ModTime)
		}
	case model.SortExtension:
		return func(i, j int) bool {
			a, b := entries[i].Name, entries[j].Name
			ra := extGroupRank(entries[i].Base, entries[i].Ext)
			rb := extGroupRank(entries[j].Base, entries[j].Ext)
			if ra != rb {
				return ra < rb
			}
			if ra == 2 {
				if compareName(entries[i].Ext, entries[j].Ext) {
					return true
				}
				if compareName(entries[j].Ext, entries[i].Ext) {
					return false
				}
			}
			return mainSort(a, b)
		}
	case model.SortNatural:
		return func(i, j int) bool {
			a, b := entries[i].Name, entries[j].Name
			if res, ok := dotFileOrder(a, b); ok {
				return res
			}
			return a < b
		}
	case model.SortNone:
		fallthrough
	default:
		return func(i, j int) bool { return i < j }
	}
}

// mainSort orders two names with dotfiles first, then case-insensitive (or
// byte-order in the C locale) within each group. "." and ".." sort as
// ordinary dotfiles.
func mainSort(a, b string) bool {
	aDot := strings.HasPrefix(a, ".")
	bDot := strings.HasPrefix(b, ".")
	if aDot != bDot {
		return aDot
	}
	if aDot {
		a = strings.TrimPrefix(a, ".")
	}
	if bDot {
		b = strings.TrimPrefix(b, ".")
	}
	return compareName(a, b)
}

// dotFileOrder reports whether the dotfile grouping alone decides the
// order; returns (result, decided).
func dotFileOrder(a, b string) (bool, bool) {
	aDot := strings.HasPrefix(a, ".")
	bDot := strings.HasPrefix(b, ".")
	if aDot != bDot {
		return aDot, true
	}
	return false, false
}

// extGroupRank ranks an entry for -X (sort by extension): 0 extensionless
// non-dotfile, 1 dotfile, 2 file-with-extension.
func extGroupRank(base, ext string) int {
	if strings.HasPrefix(base, ".") {
		return 1
	}
	if ext == "" {
		return 0
	}
	return 2
}

func compareName(a, b string) bool {
	if cLocaleCollation() {
		return a < b
	}
	return strings.ToLower(a) < strings.ToLower(b)
}

func cLocaleCollation() bool {
	loc := ""
	for _, k := range []string{"LC_ALL", "LC_COLLATE", "LANG"} {
		if v := os.Getenv(k); v != "" {
			loc = v
			break
		}
	}
	if i := strings.IndexByte(loc, '.'); i >= 0 {
		loc = loc[:i]
	}
	return loc == "" || loc == "C" || loc == "POSIX"
}
