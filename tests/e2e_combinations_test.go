package tests

import (
	"github.com/canta2899/logo-ls/internal/cli"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

func TestCombo_la(t *testing.T) {
	vfs := fakefs.New(hiddenTree())
	r := runApp(t, vfs, "-lae", "/root")
	assertGolden(t, "combo_la", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -la is long mode plus . and ..
	assertContainsLine(t, r.Stdout, `^drwxr-xr-x.*\./`)
	assertContainsLine(t, r.Stdout, `\.\./`)
}

func TestCombo_lA(t *testing.T) {
	vfs := fakefs.New(hiddenTree())
	r := runApp(t, vfs, "-lAe", "/root")
	assertGolden(t, "combo_lA_almost", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -lA is long mode with hidden files but no . or ..
	assertContains(t, r.Stdout, ".env")
	for _, l := range lines(r.Stdout) {
		if strings.HasSuffix(strings.TrimSpace(l), " ./") {
			t.Error("-lA must not include ./")
		}
	}
}

func TestCombo_lh(t *testing.T) {
	vfs := fakefs.New(fakefs.Dir("root", dirMeta("a"),
		fakefs.File("big.bin", 3*1024*1024, mtime("2026-01-01 10:00:00"), fileMeta("b")),
	))
	r := runApp(t, vfs, "-lhe", "/root")
	assertGolden(t, "combo_lh", r.Stdout)
	// Intent: -lh shows human-readable size in long mode
	assertContains(t, r.Stdout, "3M")
}

func TestCombo_lhS(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-lhSe", "/root")
	assertGolden(t, "combo_lhS", r.Stdout)
	// Intent: long + human-readable + sort by size
	assertContainsLine(t, r.Stdout, `^-rw-r--r-- .*zebra\.txt`) // largest non-dot
}

func TestCombo_ltr(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-ltre", "/root")
	assertGolden(t, "combo_ltr", r.Stdout)
	// Intent: sort by mtime, reversed -> oldest first (within sort group).
	// alpha.txt (Apr 2026) is newest, so it should be LAST in -ltr output.
	got := lines(r.Stdout)
	if len(got) == 0 || !strings.HasSuffix(got[len(got)-1], "alpha.txt") {
		t.Errorf("-ltr: expected alpha.txt last, got %v", got)
	}
}

func TestCombo_lSr(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-lSre", "/root")
	assertGolden(t, "combo_lSr", r.Stdout)
	// Intent: sort by size, reversed -> smallest non-dot first
	// (.hidden is even smaller but always sorts first).
}

func TestCombo_lah(t *testing.T) {
	vfs := fakefs.New(hiddenTree())
	r := runApp(t, vfs, "-lahe", "/root")
	assertGolden(t, "combo_lah", r.Stdout)
	// Intent: long, all, human-readable simultaneously
	assertContainsLine(t, r.Stdout, `^drwxr-xr-x.*\./`)
}

func TestCombo_li(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-lie", "/root")
	assertGolden(t, "combo_li", r.Stdout)
	// Intent: long + inode -> inode appears before mode
	assertContainsLine(t, r.Stdout, `^\s*1001 -rw-r--r--`)
}

func TestCombo_ls(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-lse", "/root")
	assertGolden(t, "combo_ls", r.Stdout)
	// Intent: long + block size -> block count appears before mode
	assertContainsLine(t, r.Stdout, `^\s*8 -rw-r--r--`)
}

func TestCombo_Rl(t *testing.T) {
	vfs := fakefs.New(deepTree())
	r := runApp(t, vfs, "-Rle", "/root")
	assertGolden(t, "combo_Rl", r.Stdout)
	// Intent: recursive long listing visits every subdir
	assertContains(t, r.Stdout, "deep.txt")
}

func TestCombo_Ra(t *testing.T) {
	vfs := fakefs.New(deepTree())
	r := runApp(t, vfs, "-Rae", "/root")
	assertGolden(t, "combo_Ra", r.Stdout)
	// Intent: recursive + all -> every visited dir lists . and ..
	// 3 subdirs (root, level1, level1/level2) each show "./" once.
	got := strings.Count(normalize(r.Stdout), "./")
	if got < 3 {
		t.Errorf("expected at least 3 './' occurrences, got %d", got)
	}
}

func TestCombo_Xr(t *testing.T) {
	vfs := fakefs.New(mixedExtTree())
	r := runApp(t, vfs, "-1Xre", "/root")
	assertGolden(t, "combo_Xr", r.Stdout)
	// Intent: sort by extension, reversed.
}

func TestCombo_Sr_short(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1Sre", "/root")
	assertGolden(t, "combo_Sr", r.Stdout)
}

func TestCombo_tr_short(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1tre", "/root")
	assertGolden(t, "combo_tr", r.Stdout)
}

func TestCombo_DefaultColumns(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-e", "/root")
	assertGolden(t, "combo_default_columns", r.Stdout)
	// Intent: default mode collapses into columns (multiple entries per line).
	got := lines(r.Stdout)
	if len(got) >= 3 {
		t.Errorf("default mode should column-pack 3 entries, got %d lines", len(got))
	}
}

func TestCombo_DirArg(t *testing.T) {
	// -d against a dir argument: prints the dir itself, not its contents.
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-de", "/root")
	assertGolden(t, "combo_d_dirarg", r.Stdout)
	if len(lines(r.Stdout)) != 1 {
		t.Errorf("expected 1 line, got %v", lines(r.Stdout))
	}
}
