package tests

import (
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/fs/fakefs"
	"github.com/canta2899/logo-ls/model"
)

func TestError_NonexistentPath(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/does-not-exist")
	assertExitCode(t, model.CodeSerious, r.ExitCode)
	if !strings.Contains(r.Stderr, "does-not-exist") {
		t.Errorf("expected error mentioning path, got stderr: %q", r.Stderr)
	}
}

func TestError_UnreadableDirectory(t *testing.T) {
	vfs := fakefs.New(fakefs.Dir("root", dirMeta("9000"),
		fakefs.File("visible.txt", 10, mtime("2026-01-01 10:00:00"), fileMeta("9001")),
		fakefs.Unreadable(fakefs.Dir("locked", dirMeta("9002"),
			fakefs.File("hidden.txt", 10, mtime("2026-01-01 10:00:00"), fileMeta("9003")),
		)),
	))
	r := runApp(t, vfs, "-1Re", "/root")
	// Recursive listing into a locked dir yields a Minor exit code.
	if r.ExitCode != model.CodeMinor {
		t.Errorf("expected CodeMinor for EACCES, got %d (stderr=%q)", r.ExitCode, r.Stderr)
	}
	assertContains(t, r.Stdout, "visible.txt")
}

func TestError_MixedValidAndInvalid(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/does-not-exist", "/root")
	assertExitCode(t, model.CodeSerious, r.ExitCode)
	// Valid arg still listed.
	assertContains(t, r.Stdout, "README.md")
	// Bad arg surfaced on stderr.
	if !strings.Contains(r.Stderr, "does-not-exist") {
		t.Errorf("expected stderr to mention bad arg, got: %q", r.Stderr)
	}
}
