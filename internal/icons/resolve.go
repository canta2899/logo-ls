package icons

import "strings"

// Resolve picks an icon for a (name, ext, indicator) triple. Rules: dirs by
// full name, files by full name, by sub-extension, then by extension, with
// hidden/exe fallbacks. Pure function — no FS access, no shared state.
func Resolve(name, ext, indicator string) *IconInfo {
	var i *IconInfo
	var ok bool

	switch indicator {
	case "/":
		i, ok = IconDir[strings.ToLower(name+ext)]
		if ok {
			break
		}
		if len(name) == 0 || name[0] == '.' {
			i = IconDef["hiddendir"]
			break
		}
		i = IconDef["dir"]
	default:
		i, ok = IconFileName[strings.ToLower(name+ext)]
		if ok {
			break
		}
		t := strings.Split(name, ".")
		if len(t) > 1 && t[0] != "" {
			i, ok = IconSubExt[strings.ToLower(t[len(t)-1]+ext)]
			if ok {
				break
			}
		}
		i, ok = IconExt[strings.ToLower(strings.TrimPrefix(ext, "."))]
		if ok {
			break
		}
		if len(name) == 0 || name[0] == '.' {
			i = IconDef["hiddenfile"]
			break
		}
		i = IconDef["file"]
	}

	if indicator == "*" && i != nil {
		if i.GetGlyph() == "" {
			i = IconDef["exe"]
		}
		i = i.AsExecutable()
	}
	if i == nil {
		i = IconDef["file"]
	}
	return i
}

// OpenDir returns the shared "open directory" icon used for headers.
func OpenDir() *IconInfo { return IconDef["diropen"] }
