package tests

import (
	"time"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

// mtime parses a "2006-01-02 15:04:05" string into a time.Time, panicking on
// error. Test fixtures use literal strings so the year is locked.
func mtime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

// fileMeta returns deterministic metadata for a regular file.
func fileMeta(inode string) fakefs.Meta {
	return fakefs.Meta{
		Owner: "alice", Group: "staff",
		Mode: 0o644, ModeStr: "-rw-r--r--",
		Inode: inode, Nlinks: 1, Blocks: 8,
	}
}

// dirMeta returns deterministic metadata for a directory.
//
// All dirs get the same fixed mtime so goldens don't drift with the year.
func dirMeta(inode string) fakefs.Meta {
	return fakefs.Meta{
		Owner: "alice", Group: "staff",
		Mode: 0o755, ModeStr: "drwxr-xr-x",
		Inode: inode, Nlinks: 2, Blocks: 0,
		Mtime: mtime("2026-01-01 00:00:00"),
	}
}

// execMeta returns deterministic metadata for an executable file.
func execMeta(inode string) fakefs.Meta {
	return fakefs.Meta{
		Owner: "alice", Group: "staff",
		Mode: 0o755, ModeStr: "-rwxr-xr-x",
		Inode: inode, Nlinks: 1, Blocks: 8,
	}
}

// smallTree: a handful of regular files plus one subdir, no hidden entries.
func smallTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("1000"),
		fakefs.File("README.md", 1234, mtime("2026-01-15 10:00:00"), fileMeta("1001")),
		fakefs.File("notes.txt", 256, mtime("2026-01-10 09:00:00"), fileMeta("1002")),
		fakefs.Dir("src", dirMeta("1003"),
			fakefs.File("main.go", 5678, mtime("2026-02-01 12:00:00"), fileMeta("1004")),
		),
	)
}

// hiddenTree: mix of dotfiles, dotdirs, regular entries.
func hiddenTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("2000"),
		fakefs.File(".env", 42, mtime("2026-01-05 08:00:00"), fileMeta("2001")),
		fakefs.File("README.md", 1234, mtime("2026-01-15 10:00:00"), fileMeta("2002")),
		fakefs.Dir(".config", dirMeta("2003"),
			fakefs.File("settings.toml", 88, mtime("2026-01-08 11:00:00"), fileMeta("2004")),
		),
		fakefs.Dir("src", dirMeta("2005"),
			fakefs.File("main.go", 256, mtime("2026-02-01 12:00:00"), fileMeta("2006")),
		),
	)
}

// sortFixture: designed to exercise every sort mode distinctly.
//
// - mixed dotfiles and regular files
// - files with no extension
// - names that sort differently in natural vs lexical (file2 vs file10)
// - deliberately distinct sizes and mtimes
func sortFixture() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("3000"),
		fakefs.File(".hidden", 1, mtime("2026-01-01 00:00:00"), fileMeta("3001")),
		fakefs.File("Makefile", 200, mtime("2026-03-01 10:00:00"), fileMeta("3002")),
		fakefs.File("README.md", 10, mtime("2026-02-01 10:00:00"), fileMeta("3003")),
		fakefs.File("file2.go", 500, mtime("2026-01-15 10:00:00"), fileMeta("3004")),
		fakefs.File("file10.go", 50, mtime("2026-01-20 10:00:00"), fileMeta("3005")),
		fakefs.File("zebra.txt", 999, mtime("2026-01-10 10:00:00"), fileMeta("3006")),
		fakefs.File("alpha.txt", 5, mtime("2026-04-01 10:00:00"), fileMeta("3007")),
	)
}

// symlinkMeta returns deterministic metadata for a symlink, with the proper
// pre-rendered mode string.
func symlinkMeta(inode string) fakefs.Meta {
	return fakefs.Meta{
		Owner: "alice", Group: "staff",
		Mode: 0o777, ModeStr: "lrwxrwxrwx",
		Inode: inode, Nlinks: 1, Blocks: 0,
		Mtime: mtime("2026-01-01 10:00:00"),
	}
}

// treeWithSymlinks: symlinks to files, to dirs, and a broken link.
func treeWithSymlinks() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("4000"),
		fakefs.File("target.txt", 100, mtime("2026-01-01 10:00:00"), fileMeta("4001")),
		fakefs.Symlink("link-file", "target.txt", symlinkMeta("4002")),
		fakefs.Dir("subdir", dirMeta("4003"),
			fakefs.File("inner.txt", 50, mtime("2026-01-02 10:00:00"), fileMeta("4004")),
		),
		fakefs.Symlink("link-dir", "subdir", symlinkMeta("4005")),
		fakefs.Symlink("link-broken", "nowhere", symlinkMeta("4006")),
	)
}

// gitRepoTree: simple repo layout, used together with WithGitStatus.
func gitRepoTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("5000"),
		fakefs.File("staged.txt", 10, mtime("2026-01-01 10:00:00"), fileMeta("5001")),
		fakefs.File("modified.txt", 20, mtime("2026-01-02 10:00:00"), fileMeta("5002")),
		fakefs.File("untracked.txt", 30, mtime("2026-01-03 10:00:00"), fileMeta("5003")),
		fakefs.File("clean.txt", 40, mtime("2026-01-04 10:00:00"), fileMeta("5004")),
	)
}

func gitRepoStatus() map[string]string {
	return map[string]string{
		"staged.txt":    "A",
		"modified.txt":  "M",
		"untracked.txt": "U",
	}
}

// deepTree: three-level fixture for recursion tests.
func deepTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("6000"),
		fakefs.File("top.txt", 10, mtime("2026-01-01 10:00:00"), fileMeta("6001")),
		fakefs.Dir("level1", dirMeta("6002"),
			fakefs.File("a.txt", 11, mtime("2026-01-02 10:00:00"), fileMeta("6003")),
			fakefs.Dir("level2", dirMeta("6004"),
				fakefs.File("deep.txt", 12, mtime("2026-01-03 10:00:00"), fileMeta("6005")),
			),
		),
	)
}

// mixedExtTree: files with same extension, no extension, and dirs — useful
// for -X (sort by extension) plus dotfile-first interaction.
func mixedExtTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("7000"),
		fakefs.File("Makefile", 50, mtime("2026-01-10 10:00:00"), fileMeta("7001")),
		fakefs.File("script", 60, mtime("2026-01-11 10:00:00"), fileMeta("7002")),
		fakefs.File("a.go", 70, mtime("2026-01-12 10:00:00"), fileMeta("7003")),
		fakefs.File("b.go", 80, mtime("2026-01-13 10:00:00"), fileMeta("7004")),
		fakefs.File("c.md", 90, mtime("2026-01-14 10:00:00"), fileMeta("7005")),
	)
}

// dotfileExtTree: dotfiles mixed with regular and extension-less files, used
// to verify that -X keeps dotfiles grouped. Expected -X order: Makefile
// (no ext), .hidden (dotfile group), a.go, README.md (by extension).
func dotfileExtTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("9000"),
		fakefs.File(".hidden", 1, mtime("2026-01-01 00:00:00"), fileMeta("9001")),
		fakefs.File("Makefile", 200, mtime("2026-03-01 10:00:00"), fileMeta("9002")),
		fakefs.File("README.md", 10, mtime("2026-02-01 10:00:00"), fileMeta("9003")),
		fakefs.File("a.go", 50, mtime("2026-01-20 10:00:00"), fileMeta("9004")),
	)
}

// dotfileGroupTree: a directory plus extensionless, dotfile (with and without
// an extension) and regular extension entries. Used to verify that -X keeps the
// whole dotfile group together (incl. a dotfile that *has* an extension, which
// must not sort among the regular extension files) and that "." / ".." are not
// force-pinned to the top.
func dotfileGroupTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("9100"),
		fakefs.Dir("src", dirMeta("9101")),
		fakefs.File("Makefile", 200, mtime("2026-03-01 10:00:00"), fileMeta("9102")),
		fakefs.File(".hidden", 1, mtime("2026-01-01 00:00:00"), fileMeta("9103")),
		fakefs.File(".config.json", 30, mtime("2026-01-02 00:00:00"), fileMeta("9104")),
		fakefs.File("main.go", 50, mtime("2026-01-20 10:00:00"), fileMeta("9105")),
		fakefs.File("app.json", 40, mtime("2026-01-21 10:00:00"), fileMeta("9106")),
	)
}

// execTree: a single executable file to exercise the '*' indicator.
func execTree() *fakefs.Entry {
	return fakefs.Dir("root", dirMeta("8000"),
		fakefs.File("regular.txt", 10, mtime("2026-01-01 10:00:00"), fileMeta("8001")),
		fakefs.File("run.sh", 20, mtime("2026-01-02 10:00:00"), execMeta("8002")),
	)
}
