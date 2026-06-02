package icons_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/canta2899/logo-ls/internal/icons"
)

// loadYAML writes data to a tempfile and loads it via LoadOverridesFromPath.
func loadYAML(t *testing.T, data string) (*icons.Override, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "icons.yaml")
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	return icons.LoadOverridesFromPath(path)
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

func TestLoadOverridesValid(t *testing.T) {
	ov, err := loadYAML(t, sampleYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ov == nil {
		t.Fatal("expected non-nil override")
	}

	got := icons.ResolveWith(ov, "main", ".rs", "")
	if got.Glyph != "X" || got.Color != [3]uint8{0x11, 0x22, 0x33} {
		t.Errorf("rs override not applied, got glyph=%q color=%v", got.Glyph, got.Color)
	}

	got = icons.ResolveWith(ov, "main", ".ts", "")
	if got.Glyph != string(rune(0xE628)) {
		t.Errorf("ts codepoint not applied, got %q", got.Glyph)
	}
	if got.Color != [3]uint8{0x11, 0xAA, 0xCC} {
		t.Errorf("short color form not expanded, got %v", got.Color)
	}

	got = icons.ResolveWith(ov, ".envrc", "", "")
	if got.Glyph != "E" {
		t.Errorf("filename override not applied, got %q", got.Glyph)
	}

	got = icons.ResolveWith(ov, "myproject", "", "/")
	if got.Glyph != "D" {
		t.Errorf("directory override not applied, got %q", got.Glyph)
	}

	got = icons.ResolveWith(ov, "lib.d", ".ts", "")
	if got.Glyph != "S" {
		t.Errorf("sub_extension override not applied, got %q", got.Glyph)
	}
}

// TestPartialOverrideColorOnly verifies that setting only color preserves the
// built-in glyph. This is the core use case for "I just want to recolor".
func TestPartialOverrideColorOnly(t *testing.T) {
	ov, err := loadYAML(t, `extensions: {go: {color: "#ff0000"}}`)
	if err != nil {
		t.Fatalf("loadYAML: %v", err)
	}
	base := icons.Resolve("main", ".go", "")
	got := icons.ResolveWith(ov, "main", ".go", "")
	if got.Glyph != base.Glyph {
		t.Errorf("expected built-in glyph %q to be preserved, got %q", base.Glyph, got.Glyph)
	}
	if got.Color != [3]uint8{0xFF, 0x00, 0x00} {
		t.Errorf("expected color override applied, got %v", got.Color)
	}
}

// TestPartialOverrideGlyphOnly verifies that setting only glyph preserves
// the built-in color.
func TestPartialOverrideGlyphOnly(t *testing.T) {
	ov, err := loadYAML(t, `extensions: {go: {glyph: "Z"}}`)
	if err != nil {
		t.Fatalf("loadYAML: %v", err)
	}
	base := icons.Resolve("main", ".go", "")
	got := icons.ResolveWith(ov, "main", ".go", "")
	if got.Glyph != "Z" {
		t.Errorf("expected glyph override applied, got %q", got.Glyph)
	}
	if got.Color != base.Color {
		t.Errorf("expected built-in color %v to be preserved, got %v", base.Color, got.Color)
	}
}

// TestPartialOverrideMissingBuiltin verifies that a partial override on a
// key with no matching built-in falls back to the default file/dir icon's
// fields for the unset side.
func TestPartialOverrideMissingBuiltin(t *testing.T) {
	ov, err := loadYAML(t, `extensions: {totallynew: {color: "#abcdef"}}`)
	if err != nil {
		t.Fatalf("loadYAML: %v", err)
	}
	got := icons.ResolveWith(ov, "foo", ".totallynew", "")
	// No built-in for "totallynew", so the default file icon's glyph is kept.
	def := icons.Resolve("foo", ".totallynew", "")
	if got.Glyph != def.Glyph {
		t.Errorf("expected default-file glyph %q kept, got %q", def.Glyph, got.Glyph)
	}
	if got.Color != [3]uint8{0xAB, 0xCD, 0xEF} {
		t.Errorf("expected color override applied, got %v", got.Color)
	}
}

func TestResolveWithFallsBackToDefault(t *testing.T) {
	ov, err := loadYAML(t, sampleYAML)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	withOv := icons.ResolveWith(ov, "main", ".cpp", "")
	without := icons.Resolve("main", ".cpp", "")
	if withOv.Glyph != without.Glyph || withOv.Color != without.Color {
		t.Errorf("expected built-in icon for unmatched ext, got override")
	}
}

func TestResolveWithNilOverride(t *testing.T) {
	got := icons.ResolveWith(nil, "main", ".go", "")
	want := icons.Resolve("main", ".go", "")
	if got != want {
		t.Errorf("nil override should match Resolve")
	}
}

func TestOverrideExecutableHandling(t *testing.T) {
	ov, err := loadYAML(t, `extensions:
  sh:
    glyph: "S"
    color: "#ff00ff"
`)
	if err != nil {
		t.Fatal(err)
	}
	got := icons.ResolveWith(ov, "run", ".sh", "*")
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
			ov, err := loadYAML(t, c.yaml)
			if err != nil {
				t.Fatalf("loadYAML: %v", err)
			}
			got := icons.ResolveWith(ov, "main", ".rs", "").GetGlyph()
			if got != want {
				t.Errorf("got glyph %q, want %q", got, want)
			}
		})
	}
}

func TestGlyphLiteralIsNotMisparsed(t *testing.T) {
	ov, err := loadYAML(t, `extensions: {rs: {glyph: "abc", color: "#000000"}}`)
	if err != nil {
		t.Fatal(err)
	}
	if got := icons.ResolveWith(ov, "main", ".rs", "").GetGlyph(); got != "abc" {
		t.Errorf("expected literal %q, got %q", "abc", got)
	}
}

func TestLoadOverridesEmptyFile(t *testing.T) {
	ov, err := loadYAML(t, "   \n")
	if err != nil || ov != nil {
		t.Errorf("expected (nil, nil), got (%v, %v)", ov, err)
	}
}

func TestLoadOverridesInvalidColor(t *testing.T) {
	_, err := loadYAML(t, `extensions:
  rs:
    glyph: "X"
    color: "not-a-color"
`)
	if err == nil {
		t.Error("expected error on bad color")
	}
}

// TestLoadOverridesEmptyEntry verifies that an entry with neither glyph nor
// color is rejected (every entry must do at least one thing).
func TestLoadOverridesEmptyEntry(t *testing.T) {
	_, err := loadYAML(t, `extensions:
  rs: {}
`)
	if err == nil {
		t.Error("expected error when entry has no glyph and no color")
	}
}

func TestLoadOverridesMalformedYAML(t *testing.T) {
	_, err := loadYAML(t, "extensions: [: oops")
	if err == nil {
		t.Error("expected error on malformed yaml")
	}
}

func TestLoadOverridesDiscovery(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	path := filepath.Join(home, ".logo-ls-overrides.yaml")
	if err := os.WriteFile(path, []byte(sampleYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	ov, err := icons.LoadOverrides()
	if err != nil {
		t.Fatalf("LoadOverrides: %v", err)
	}
	if ov == nil || ov.Source != path {
		t.Fatalf("expected override loaded from %s, got %+v", path, ov)
	}
}

func TestLoadOverridesDiscoveryXDG(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", xdg)

	dir := filepath.Join(xdg, "logo-ls")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "logo-ls-overrides.yaml")
	if err := os.WriteFile(path, []byte(sampleYAML), 0o644); err != nil {
		t.Fatal(err)
	}

	ov, err := icons.LoadOverrides()
	if err != nil {
		t.Fatalf("LoadOverrides: %v", err)
	}
	if ov == nil || ov.Source != path {
		t.Fatalf("expected override loaded from XDG path %s, got %+v", path, ov)
	}
}

func TestLoadOverridesMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	ov, err := icons.LoadOverrides()
	if err != nil {
		t.Fatalf("expected nil error when no file, got %v", err)
	}
	if ov != nil {
		t.Fatalf("expected nil override, got %+v", ov)
	}
}
