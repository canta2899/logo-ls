package tests

import (
	"github.com/canta2899/logo-ls/internal/cli"
	"testing"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

// Symlink tests: file-link, dir-link, and broken link must render with the
// '@' indicator (or in long mode, the 'l' mode prefix). The fakefs symlink
// targets are bare names resolved relative to the parent dir.

func TestSymlink_OneLineIndicator(t *testing.T) {
	vfs := fakefs.New(treeWithSymlinks())
	r := runApp(t, vfs, "-1e", "/root")
	assertGolden(t, "symlink_oneline", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// All three links rendered with '@'.
	assertContains(t, r.Stdout, "link-file@")
	assertContains(t, r.Stdout, "link-dir@")
	assertContains(t, r.Stdout, "link-broken@")
}

func TestSymlink_LongModeShowsLPrefix(t *testing.T) {
	vfs := fakefs.New(treeWithSymlinks())
	r := runApp(t, vfs, "-le", "/root")
	assertGolden(t, "symlink_long", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Long mode: mode column starts with 'l', resolvable links show "~> target".
	assertContainsLine(t, r.Stdout, `^l.*link-file ~> /root/target\.txt`)
	assertContainsLine(t, r.Stdout, `^l.*link-dir ~> /root/subdir`)
	assertContainsLine(t, r.Stdout, `^l.*link-broken`)
}
