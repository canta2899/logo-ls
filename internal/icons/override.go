package icons

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Override holds user-supplied icon overrides loaded from a YAML file.
// A nil *Override is valid and corresponds to "no overrides".
type Override struct {
	Source string // absolute path of the YAML file

	dirs      map[string]overrideEntry
	fileNames map[string]overrideEntry
	subExts   map[string]overrideEntry
	exts      map[string]overrideEntry
}

// overrideEntry is the parsed form of a single YAML entry.
type overrideEntry struct {
	glyph    string
	hasColor bool
	color    [3]uint8
}

type yamlEntry struct {
	Glyph     string `yaml:"glyph"`
	Codepoint string `yaml:"codepoint"`
	Color     string `yaml:"color"`
}

type yamlFile struct {
	Directories   map[string]yamlEntry `yaml:"directories"`
	Files         map[string]yamlEntry `yaml:"files"`
	Extensions    map[string]yamlEntry `yaml:"extensions"`
	SubExtensions map[string]yamlEntry `yaml:"sub_extensions"`
}

// candidatePaths returns the YAML file locations checked at startup, in order.
func candidatePaths() []string {
	var out []string
	home, _ := os.UserHomeDir()
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		out = append(out, filepath.Join(xdg, "logo-ls", "logo-ls-overrides.yaml"))
	} else if home != "" {
		out = append(out, filepath.Join(home, ".config", "logo-ls", "logo-ls-overrides.yaml"))
	}
	if home != "" {
		out = append(out, filepath.Join(home, ".logo-ls-overrides.yaml"))
	}
	return out
}

// LoadOverrides discovers and parses the user's override file. Returns
// (nil, nil) when no file is found or when the file is empty.
func LoadOverrides() (*Override, error) {
	for _, path := range candidatePaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		return parseOverrides(path, data)
	}
	return nil, nil
}

// LoadOverridesFromPath reads and parses an explicit file path, bypassing
// the default discovery locations.
func LoadOverridesFromPath(path string) (*Override, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return parseOverrides(path, data)
}

func parseOverrides(path string, data []byte) (*Override, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var raw yamlFile
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	ov := &Override{Source: path}
	var err error
	if ov.dirs, err = buildEntryMap(raw.Directories); err != nil {
		return nil, fmt.Errorf("%s directories: %w", path, err)
	}
	if ov.fileNames, err = buildEntryMap(raw.Files); err != nil {
		return nil, fmt.Errorf("%s files: %w", path, err)
	}
	if ov.exts, err = buildEntryMap(raw.Extensions); err != nil {
		return nil, fmt.Errorf("%s extensions: %w", path, err)
	}
	if ov.subExts, err = buildEntryMap(raw.SubExtensions); err != nil {
		return nil, fmt.Errorf("%s sub_extensions: %w", path, err)
	}
	if ov.empty() {
		return nil, nil
	}
	return ov, nil
}

func (o *Override) empty() bool {
	return len(o.dirs)+len(o.fileNames)+len(o.subExts)+len(o.exts) == 0
}

func buildEntryMap(in map[string]yamlEntry) (map[string]overrideEntry, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make(map[string]overrideEntry, len(in))
	for key, entry := range in {
		e, err := parseEntry(entry)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", key, err)
		}
		out[strings.ToLower(strings.TrimPrefix(key, "."))] = e
		out[strings.ToLower(key)] = e
	}
	return out, nil
}

func parseEntry(e yamlEntry) (overrideEntry, error) {
	out := overrideEntry{}
	if e.Glyph != "" || e.Codepoint != "" {
		g, err := parseGlyph(e.Glyph, e.Codepoint)
		if err != nil {
			return out, err
		}
		out.glyph = g
	}
	if e.Color != "" {
		c, err := parseColor(e.Color)
		if err != nil {
			return out, err
		}
		out.hasColor = true
		out.color = c
	}
	if out.glyph == "" && !out.hasColor {
		return out, errors.New("override must set at least one of glyph, codepoint, or color")
	}
	return out, nil
}

func parseGlyph(literal, codepoint string) (string, error) {
	if codepoint != "" {
		return parseCodepoint(codepoint)
	}
	// The glyph field also accepts codepoint syntax ("U+e7a8", "0xe7a8")
	// so users don't have to remember YAML's "\uXXXX" escape rules.
	if hasCodepointPrefix(literal) {
		return parseCodepoint(literal)
	}
	return literal, nil
}

func hasCodepointPrefix(s string) bool {
	s = strings.TrimSpace(s)
	switch {
	case strings.HasPrefix(s, "U+"), strings.HasPrefix(s, "u+"):
		return true
	case strings.HasPrefix(s, "0x"), strings.HasPrefix(s, "0X"):
		return true
	}
	return false
}

func parseCodepoint(s string) (string, error) {
	raw := s
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "U+")
	s = strings.TrimPrefix(s, "u+")
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return "", fmt.Errorf("invalid codepoint %q: %w", raw, err)
	}
	return string(rune(n)), nil
}

func parseColor(s string) ([3]uint8, error) {
	var out [3]uint8
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 3:
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 6:
	default:
		return out, fmt.Errorf("invalid color %q: expected #RGB or #RRGGBB", s)
	}
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return out, fmt.Errorf("invalid color %q: %w", s, err)
	}
	out[0] = uint8(n >> 16)
	out[1] = uint8(n >> 8)
	out[2] = uint8(n)
	return out, nil
}

// lookupEntry returns the first matching override entry across the lookup
func (o *Override) lookupEntry(name, fileExt, indicator string) (overrideEntry, bool) {
	if o == nil {
		return overrideEntry{}, false
	}
	lowerNameExt := strings.ToLower(name + fileExt)
	if indicator == "/" {
		e, ok := o.dirs[lowerNameExt]
		return e, ok
	}
	if e, ok := o.fileNames[lowerNameExt]; ok {
		return e, true
	}
	t := strings.Split(name, ".")
	if len(t) > 1 && t[0] != "" {
		subKey := strings.ToLower(t[len(t)-1] + fileExt)
		if e, ok := o.subExts[subKey]; ok {
			return e, true
		}
	}
	extKey := strings.ToLower(strings.TrimPrefix(fileExt, "."))
	if e, ok := o.exts[extKey]; ok {
		return e, true
	}
	return overrideEntry{}, false
}

// apply returns a copy of base with the entry's set fields swapped in.
func (e overrideEntry) apply(base *IconInfo) *IconInfo {
	if base == nil {
		out := &IconInfo{Glyph: e.glyph, Color: e.color}
		return out
	}
	out := *base
	if e.glyph != "" {
		out.Glyph = e.glyph
	}
	if e.hasColor {
		out.Color = e.color
	}
	return &out
}
