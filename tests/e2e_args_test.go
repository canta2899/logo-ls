package tests

import (
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/fs/fakefs"
	"github.com/canta2899/logo-ls/model"
)

func TestArgs_FileAndDirectory(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/root/README.md", "/root")
	assertGolden(t, "args_file_and_dir", r.Stdout)
	assertExitCode(t, model.CodeOk, r.ExitCode)
	// File appears first (before directory listing).
	assertContains(t, r.Stdout, "README.md")
	assertContains(t, r.Stdout, "notes.txt")
}

func twoDirRoot() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("100"),
		fakefs.Dir("alpha", dirMeta("101"),
			fakefs.File("a.txt", 10, mtime("2026-01-01 10:00:00"), fileMeta("102")),
		),
		fakefs.Dir("beta", dirMeta("103"),
			fakefs.File("b.txt", 20, mtime("2026-01-02 10:00:00"), fileMeta("104")),
		),
	)
}

func TestArgs_TwoDirectories(t *testing.T) {
	vfs := fakefs.New(twoDirRoot())
	r := runApp(t, vfs, "-1e", "/root/alpha", "/root/beta")
	assertGolden(t, "args_two_dirs", r.Stdout)
	assertExitCode(t, model.CodeOk, r.ExitCode)
	assertContains(t, r.Stdout, "/root/alpha:")
	assertContains(t, r.Stdout, "/root/beta:")
	assertContains(t, r.Stdout, "a.txt")
	assertContains(t, r.Stdout, "b.txt")
}

func TestArgs_DirOrderingIsSorted(t *testing.T) {
	// FileList is sorted before processing, so passing args in reverse should
	// produce the same output as passing them in sorted order.
	sorted := runApp(t, fakefs.New(twoDirRoot()), "-1e", "/root/alpha", "/root/beta").Stdout
	reversed := runApp(t, fakefs.New(twoDirRoot()), "-1e", "/root/beta", "/root/alpha").Stdout
	if normalize(sorted) != normalize(reversed) {
		t.Errorf("dir ordering not normalized:\nsorted=%q\nreversed=%q", sorted, reversed)
	}
	if strings.Index(sorted, "/root/alpha:") >= strings.Index(sorted, "/root/beta:") {
		t.Errorf("expected /root/alpha before /root/beta, got:\n%s", sorted)
	}
}
