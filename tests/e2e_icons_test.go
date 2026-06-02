package tests

import (
	"github.com/canta2899/logo-ls/internal/cli"
	"testing"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

// Icon tests (behavioral-only by  asserting that default output emits nerd-font glyphs)

func inPUA(s string) bool {
	for _, r := range s {
		if (r >= 0xE000 && r <= 0xF8FF) ||
			(r >= 0xF0000 && r <= 0xFFFFD) ||
			(r >= 0x100000 && r <= 0x10FFFD) {
			return true
		}
	}
	return false
}

func TestIcons_DefaultEmitsGlyphs(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "/root")
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	if !inPUA(r.Stdout) {
		t.Errorf("expected nerd-font glyph in default output, got:\n%s", r.Stdout)
	}
}

func TestIcons_FlagESuppressesGlyphs(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-e", "/root")
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	if inPUA(r.Stdout) {
		t.Errorf("expected no nerd-font glyph with -e, got:\n%s", r.Stdout)
	}
}

func TestIcons_LongModeStillHasGlyphs(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-l", "/root")
	if !inPUA(r.Stdout) {
		t.Errorf("expected glyphs in -l output, got:\n%s", r.Stdout)
	}
}

func TestIcons_HiddenFileHasGlyph(t *testing.T) {
	vfs := fakefs.New(hiddenTree())
	r := runApp(t, vfs, "-A", "/root")
	if !inPUA(r.Stdout) {
		t.Errorf("expected glyph for hidden entries, got:\n%s", r.Stdout)
	}
}
