package icons

import "strings"

// Resolve picks an icon for a (name, ext, indicator) triple.
func Resolve(name, ext, indicator string) *IconInfo {
	return ResolveWith(nil, name, ext, indicator)
}

// ResolveWith is like Resolve but applies a user-provided Override on top of
// the built-in lookup.
func ResolveWith(ov *Override, name, fileExt, indicator string) *IconInfo {
	base := resolveBuiltin(name, fileExt, indicator)
	if entry, ok := ov.lookupEntry(name, fileExt, indicator); ok {
		base = entry.apply(base)
	}
	return applyIndicator(base, indicator)
}

// resolveBuiltin runs the built-in lookup chain without any override
func resolveBuiltin(name, fileExt, indicator string) *IconInfo {
	lowerNameExt := strings.ToLower(name + fileExt)

	if indicator == "/" {
		if i, ok := IconDir[lowerNameExt]; ok {
			return i
		}
		if len(name) == 0 || name[0] == '.' {
			return IconDef["hiddendir"]
		}
		return IconDef["dir"]
	}

	if i, ok := IconFileName[lowerNameExt]; ok {
		return i
	}
	t := strings.Split(name, ".")
	if len(t) > 1 && t[0] != "" {
		subKey := strings.ToLower(t[len(t)-1] + fileExt)
		if i, ok := IconSubExt[subKey]; ok {
			return i
		}
	}
	extKey := strings.ToLower(strings.TrimPrefix(fileExt, "."))
	if i, ok := IconExt[extKey]; ok {
		return i
	}
	if len(name) == 0 || name[0] == '.' {
		return IconDef["hiddenfile"]
	}
	return IconDef["file"]
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
