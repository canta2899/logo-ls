package icons

import "strings"

// Resolve picks an icon for a (name, ext, indicator) triple.
func Resolve(name, ext, indicator string) *IconInfo {
	return ResolveWith(nil, name, ext, indicator)
}

// ResolveWith is like Resolve but consults a user-provided Extension first
func ResolveWith(ext *Extension, name, fileExt, indicator string) *IconInfo {
	var i *IconInfo
	var ok bool

	lowerNameExt := strings.ToLower(name + fileExt)

	switch indicator {
	case "/":
		if i = ext.lookupDir(lowerNameExt); i != nil {
			break
		}
		i, ok = IconDir[lowerNameExt]
		if ok {
			break
		}
		if len(name) == 0 || name[0] == '.' {
			i = IconDef["hiddendir"]
			break
		}
		i = IconDef["dir"]
	default:
		if i = ext.lookupFileName(lowerNameExt); i != nil {
			break
		}
		i, ok = IconFileName[lowerNameExt]
		if ok {
			break
		}
		t := strings.Split(name, ".")
		if len(t) > 1 && t[0] != "" {
			subKey := strings.ToLower(t[len(t)-1] + fileExt)
			if i = ext.lookupSubExt(subKey); i != nil {
				break
			}
			i, ok = IconSubExt[subKey]
			if ok {
				break
			}
		}
		extKey := strings.ToLower(strings.TrimPrefix(fileExt, "."))
		if i = ext.lookupExt(extKey); i != nil {
			break
		}
		i, ok = IconExt[extKey]
		if ok {
			break
		}
		if len(name) == 0 || name[0] == '.' {
			i = IconDef["hiddenfile"]
			break
		}
		i = IconDef["file"]
	}

	return applyIndicator(i, indicator)
}

// applyIndicator handles the "*" executable special case and the nil fallback.
func applyIndicator(i *IconInfo, indicator string) *IconInfo {
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
