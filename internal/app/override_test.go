package app

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/internal/icons"
)

// TestAppIconOverride verifies that an override wired into the app flows
// through to the inspector and overrides the resolved icon for matching
// files while leaving non-matching files on the default resolver.
func TestAppIconOverride(t *testing.T) {
	tempDir := t.TempDir()
	hit := filepath.Join(tempDir, "main.rs")
	miss := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(hit, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(miss, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	yamlPath := filepath.Join(t.TempDir(), "icons.yaml")
	if err := os.WriteFile(yamlPath, []byte(`
extensions:
  rs:
    glyph: "OVR"
    color: "#010203"
`), 0o644); err != nil {
		t.Fatal(err)
	}
	ov, err := icons.LoadOverridesFromPath(yamlPath)
	if err != nil {
		t.Fatalf("LoadOverridesFromPath: %v", err)
	}

	conf := &cli.Config{
		AllMode:         cli.IncludeDefault,
		LongListingMode: cli.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
	}
	appInstance := newTestApp(conf, log.New(io.Discard, "", 0), new(bytes.Buffer))
	appInstance.IconOverride = ov

	f, err := appInstance.FS.Open(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	dirEntry := &DirectoryEntry{File: f, AbsPath: tempDir}

	d, err := appInstance.ProcessDirectory(dirEntry)
	if err != nil {
		t.Fatalf("ProcessDirectory: %v", err)
	}

	var rsIcon, goIcon string
	for _, e := range d.Files {
		switch e.Name {
		case "main.rs":
			rsIcon = e.Icon.GetGlyph()
		case "main.go":
			goIcon = e.Icon.GetGlyph()
		}
	}
	if rsIcon != "OVR" {
		t.Errorf("expected rs icon overridden, got %q", rsIcon)
	}
	if goIcon == "" || goIcon == "OVR" {
		t.Errorf("expected default go icon, got %q", goIcon)
	}
}

func TestAppNoIconOverrideMatchesDefault(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "main.rs")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	conf := &cli.Config{
		AllMode:         cli.IncludeDefault,
		LongListingMode: cli.LongListingNone,
		TimeFormatter:   DummyTimeFormatter{},
	}
	appInstance := newTestApp(conf, log.New(io.Discard, "", 0), new(bytes.Buffer))

	f, err := appInstance.FS.Open(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	d, err := appInstance.ProcessDirectory(&DirectoryEntry{File: f, AbsPath: tempDir})
	if err != nil {
		t.Fatalf("ProcessDirectory: %v", err)
	}

	want := icons.Resolve("main", ".rs", "").GetGlyph()
	for _, e := range d.Files {
		if e.Name == "main.rs" {
			if got := e.Icon.GetGlyph(); got != want {
				t.Errorf("nil override should match default Resolve, got %q want %q", got, want)
			}
		}
	}
}
