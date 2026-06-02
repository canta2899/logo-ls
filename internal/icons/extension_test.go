package icons_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canta2899/logo-ls/internal/icons"
)

// loadYAML writes data to a tempfile and loads it via LoadExtensionFromPath.
func loadYAML(t *testing.T, data string) (*icons.Extension, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "icons.yaml")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return icons.LoadExtensionFromPath(path)
}

const sampleYAML = `
extensions:
  rs:
    glyph: "X"
    color: "#112233"
  ts:
    codepoint: "U+E628"
    color: "1ac"
files:
  .envrc:
    glyph: "E"
    color: "#abcdef"
directories:
  myproject:
    glyph: "D"
    color: "#000000"
sub_extensions:
  d.ts:
    glyph: "S"
    color: "#ffffff"
`

func TestLoadExtensionFromValid(t *testing.T) {
	ext, err := loadYAML(t, sampleYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ext == nil {
		t.Fatal("expected non-nil extension")
	}

	got := icons.ResolveWith(ext, "main", ".rs", "")
	if got.Glyph != "X" || got.Color != [3]uint8{0x11, 0x22, 0x33} {
		t.Errorf("rs override not applied, got glyph=%q color=%v", got.Glyph, got.Color)
	}

	got = icons.ResolveWith(ext, "main", ".ts", "")
	if got.Glyph != string(rune(0xE628)) {
		t.Errorf("ts codepoint not applied, got %q", got.Glyph)
	}
	if got.Color != [3]uint8{0x11, 0xAA, 0xCC} {
		t.Errorf("short color form not expanded, got %v", got.Color)
	}

	got = icons.ResolveWith(ext, ".envrc", "", "")
	if got.Glyph != "E" {
		t.Errorf("filename override not applied, got %q", got.Glyph)
	}

	got = icons.ResolveWith(ext, "myproject", "", "/")
	if got.Glyph != "D" {
		t.Errorf("directory override not applied, got %q", got.Glyph)
	}

	got = icons.ResolveWith(ext, "lib.d", ".ts", "")
	if got.Glyph != "S" {
		t.Errorf("sub_extension override not applied, got %q", got.Glyph)
	}
}

func TestResolveWithFallsBackToDefault(t *testing.T) {
	ext, err := loadYAML(t, sampleYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	withExt := icons.ResolveWith(ext, "main", ".go", "")
	without := icons.Resolve("main", ".go", "")
	if withExt != without {
		t.Errorf("expected fallback to built-in icon, got override")
	}
}

func TestResolveWithNilExtension(t *testing.T) {
	got := icons.ResolveWith(nil, "main", ".go", "")
	want := icons.Resolve("main", ".go", "")
	if got != want {
		t.Errorf("nil extension should match Resolve")
	}
}

func TestExtensionExecutableHandling(t *testing.T) {
	ext, err := loadYAML(t, `extensions:
  sh:
    glyph: "S"
    color: "#ff00ff"
`)
	if err != nil {
		t.Fatal(err)
	}
	got := icons.ResolveWith(ext, "run", ".sh", "*")
	if got.Glyph != "S" {
		t.Errorf("override glyph lost for executable, got %q", got.Glyph)
	}
	if !got.IsExecutable {
		t.Errorf("expected IsExecutable=true for indicator=*")
	}
}

func TestGlyphAcceptsCodepointSyntax(t *testing.T) {
	cases := []struct {
		name string
		yaml string
	}{
		{"U+ form", `extensions: {rs: {glyph: "U+E7A8", color: "#000000"}}`},
		{"u+ form", `extensions: {rs: {glyph: "u+e7a8", color: "#000000"}}`},
		{"0x form", `extensions: {rs: {glyph: "0xE7A8", color: "#000000"}}`},
		{"codepoint field", `extensions: {rs: {codepoint: "E7A8", color: "#000000"}}`},
	}
	want := string(rune(0xE7A8))
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ext, err := loadYAML(t, c.yaml)
			if err != nil {
				t.Fatalf("loadYAML: %v", err)
			}
			got := icons.ResolveWith(ext, "main", ".rs", "").GetGlyph()
			if got != want {
				t.Errorf("got glyph %q, want %q", got, want)
			}
		})
	}
}

func TestGlyphLiteralIsNotMisparsed(t *testing.T) {
	// A short alphanumeric string with no U+/0x prefix must be treated as a
	// literal, not parsed as hex — even when it happens to be valid hex.
	ext, err := loadYAML(t, `extensions: {rs: {glyph: "abc", color: "#000000"}}`)
	if err != nil {
		t.Fatal(err)
	}
	if got := icons.ResolveWith(ext, "main", ".rs", "").GetGlyph(); got != "abc" {
		t.Errorf("expected literal %q, got %q", "abc", got)
	}
}

func TestLoadExtensionEmptyFile(t *testing.T) {
	ext, err := loadYAML(t, "   \n")
	if err != nil || ext != nil {
		t.Errorf("expected (nil, nil), got (%v, %v)", ext, err)
	}
}

func TestLoadExtensionInvalidColor(t *testing.T) {
	_, err := loadYAML(t, `extensions:
  rs:
    glyph: "X"
    color: "not-a-color"
`)
	if err == nil {
		t.Error("expected error on bad color")
	}
}

func TestLoadExtensionMissingGlyph(t *testing.T) {
	_, err := loadYAML(t, `extensions:
  rs:
    color: "#112233"
`)
	if err == nil {
		t.Error("expected error when glyph and codepoint are both missing")
	}
}

func TestLoadExtensionMalformedYAML(t *testing.T) {
	_, err := loadYAML(t, "extensions: [: oops")
	if err == nil {
		t.Error("expected error on malformed yaml")
	}
}

// TestLoadExtensionDiscovery verifies the path lookup logic by pointing
// HOME at a tempdir holding the user file.
func TestLoadExtensionDiscovery(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	path := filepath.Join(home, ".logo-ls-icons.yaml")
	if err := os.WriteFile(path, []byte(sampleYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	ext, err := icons.LoadExtension()
	if err != nil {
		t.Fatalf("LoadExtension: %v", err)
	}
	if ext == nil || ext.Source != path {
		t.Fatalf("expected extension loaded from %s, got %+v", path, ext)
	}
}

func TestLoadExtensionDiscoveryXDG(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", xdg)

	dir := filepath.Join(xdg, "logo-ls")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "logo-ls-icons.yaml")
	if err := os.WriteFile(path, []byte(sampleYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	ext, err := icons.LoadExtension()
	if err != nil {
		t.Fatalf("LoadExtension: %v", err)
	}
	if ext == nil || ext.Source != path {
		t.Fatalf("expected extension loaded from XDG path %s, got %+v", path, ext)
	}
}

func TestLoadExtensionMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	ext, err := icons.LoadExtension()
	if err != nil {
		t.Fatalf("expected nil error when no file, got %v", err)
	}
	if ext != nil {
		t.Fatalf("expected nil extension, got %+v", ext)
	}
}
