// Package format owns small formatting and icon-resolution helpers shared
// by the renderer and the inspector.
package format

import (
	"fmt"
	"strings"

	"github.com/canta2899/logo-ls/icons"
)

// GetOpenDirIcon returns the shared "open directory" icon used for headers.
func GetOpenDirIcon() *icons.IconInfo {
	return icons.IconDef["diropen"]
}

// GetIcon picks the icon for a (name, ext, indicator) triple. The rules
// mirror the original ls-icons behaviour: dirs by name, files by full
// name, by sub-extension, then by extension, with hidden/exe fallbacks.
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

	if indicator == "*" && i != nil {
		if i.GetGlyph() == "" {
			i = icons.IconDef["exe"]
		}
		i = i.AsExecutable()
	}

	if i == nil {
		i = icons.IconDef["file"]
	}

	return i
}

// GetFormattedSize renders b as either a raw byte count or a human-readable
// 1K/2.3M-style size when humanReadable is true.
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
	if size == float64(int64(size)) {
		return fmt.Sprintf("%d%c", int64(size), "KMGTPE"[exp])
	}
	return fmt.Sprintf("%.1f%c", size, "KMGTPE"[exp])
}
