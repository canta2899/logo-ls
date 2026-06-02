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

// Extension holds user-supplied icon overrides loaded from a YAML file.
// A nil *Extension is valid and corresponds to "no overrides"
type Extension struct {
	Source string // absolute path of the YAML file

	dirs      map[string]*IconInfo
	fileNames map[string]*IconInfo
	subExts   map[string]*IconInfo
	exts      map[string]*IconInfo
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

// looks for configuration files in the classic order
func candidatePaths() []string {
	var out []string
	home, _ := os.UserHomeDir()
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		out = append(out, filepath.Join(xdg, "logo-ls", "logo-ls-icons.yaml"))
	} else if home != "" {
		out = append(out, filepath.Join(home, ".config", "logo-ls", "logo-ls-icons.yaml"))
	}
	if home != "" {
		out = append(out, filepath.Join(home, ".logo-ls-icons.yaml"))
	}
	return out
}

// LoadExtension discovers and parses the user's icon override file. Returns
// (nil, nil) when no file is found or when the file is empty. Returns an error
// only when a file exists but cannot be parsed or contains invalid entries.
func LoadExtension() (*Extension, error) {
	for _, path := range candidatePaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		return parseExtension(path, data)
	}
	return nil, nil
}

// LoadExtensionFromPath reads and parses an explicit file path, bypassing the
// default discovery locations.
func LoadExtensionFromPath(path string) (*Extension, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return parseExtension(path, data)
}

func parseExtension(path string, data []byte) (*Extension, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var raw yamlFile
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	ext := &Extension{Source: path}
	var err error
	if ext.dirs, err = buildMap(raw.Directories); err != nil {
		return nil, fmt.Errorf("%s directories: %w", path, err)
	}
	if ext.fileNames, err = buildMap(raw.Files); err != nil {
		return nil, fmt.Errorf("%s files: %w", path, err)
	}
	if ext.exts, err = buildMap(raw.Extensions); err != nil {
		return nil, fmt.Errorf("%s extensions: %w", path, err)
	}
	if ext.subExts, err = buildMap(raw.SubExtensions); err != nil {
		return nil, fmt.Errorf("%s sub_extensions: %w", path, err)
	}
	if ext.empty() {
		return nil, nil
	}
	return ext, nil
}

func (e *Extension) empty() bool {
	return len(e.dirs)+len(e.fileNames)+len(e.subExts)+len(e.exts) == 0
}

func buildMap(in map[string]yamlEntry) (map[string]*IconInfo, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make(map[string]*IconInfo, len(in))
	for key, entry := range in {
		info, err := entryToIcon(entry)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", key, err)
		}
		// Keys are matched case-insensitively, so normalize here once.
		out[strings.ToLower(strings.TrimPrefix(key, "."))] = info
		// Also store the raw key for filename/dir lookups where dots are kept.
		out[strings.ToLower(key)] = info
	}
	return out, nil
}

func entryToIcon(e yamlEntry) (*IconInfo, error) {
	glyph, err := parseGlyph(e.Glyph, e.Codepoint)
	if err != nil {
		return nil, err
	}
	color, err := parseColor(e.Color)
	if err != nil {
		return nil, err
	}
	return &IconInfo{Glyph: glyph, Color: color}, nil
}

func parseGlyph(literal, codepoint string) (string, error) {
	// codepoint field always parses as hex.
	if codepoint != "" {
		return parseCodepoint(codepoint)
	}
	if literal == "" {
		return "", errors.New("missing glyph or codepoint")
	}
	// The glyph field also accepts codepoint syntax ("U+e7a8", "0xe7a8")
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
	if s == "" {
		return out, errors.New("missing color")
	}
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 3:
		// Expand "rgb" -> "rrggbb".
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

func (e *Extension) lookupDir(key string) *IconInfo {
	if e == nil {
		return nil
	}
	return e.dirs[key]
}

func (e *Extension) lookupFileName(key string) *IconInfo {
	if e == nil {
		return nil
	}
	return e.fileNames[key]
}

func (e *Extension) lookupSubExt(key string) *IconInfo {
	if e == nil {
		return nil
	}
	return e.subExts[key]
}

func (e *Extension) lookupExt(key string) *IconInfo {
	if e == nil {
		return nil
	}
	return e.exts[key]
}
