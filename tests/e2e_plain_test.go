package tests

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/internal/app"
	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

// runPlain runs the app with the non-tty settings from main and returns the
// raw stdout. runApp normalizes ANSI codes and padding away, so it can't be
// used to check them.
func runPlain(t *testing.T, root *fakefs.Entry, flag string) string {
	t.Helper()
	cfg, _, err := cli.BuildConfig([]string{"logo-ls", "--no-override", flag, "/root"})
	if err != nil {
		t.Fatal(err)
	}
	cfg.TimeFormatter = fixedTime{}
	cfg.NoColor = true
	cfg.NoIndicator = true

	var stdout bytes.Buffer
	a := &app.App{
		Config:   cfg,
		Writer:   &stdout,
		Logger:   log.New(&bytes.Buffer{}, "", 0),
		FS:       fakefs.New(root),
		ExitCode: cli.CodeOk,
	}
	a.Run()
	return stdout.String()
}

// Non-tty output must be plain: no ANSI codes and no trailing padding.
func TestPlainOutput_NoANSINoTrailingSpaces(t *testing.T) {
	for _, flag := range []string{"-1", "-l"} {
		out := runPlain(t, smallTree(), flag)
		if strings.Contains(out, "\x1b") {
			t.Errorf("%s: unexpected ANSI codes: %q", flag, out)
		}
		for _, l := range strings.Split(out, "\n") {
			if l != strings.TrimRight(l, " ") {
				t.Errorf("%s: trailing spaces in %q", flag, l)
			}
		}
	}
}

// Non-tty output must print real file names, without /, *, @ suffixes.
func TestPlainOutput_NoIndicators(t *testing.T) {
	out := runPlain(t, fakefs.Dir("root", dirMeta("9000"),
		fakefs.File("run.sh", 20, mtime("2026-01-02 10:00:00"), execMeta("9001")),
		fakefs.Dir("subdir", dirMeta("9002")),
		fakefs.Symlink("link", "subdir", symlinkMeta("9003")),
	), "-1e")
	want := "link\nrun.sh\nsubdir\n"
	if out != want {
		t.Errorf("want %q, got %q", want, out)
	}
}

// Long mode keeps the symlink target, like ls -l does.
func TestPlainOutput_LongKeepsSymlinkTarget(t *testing.T) {
	out := runPlain(t, treeWithSymlinks(), "-le")
	assertContainsLine(t, out, `link-dir ~> /root/subdir$`)
	assertContainsLine(t, out, ` subdir$`)
}
