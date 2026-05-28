package fakefs

import (
	"io/fs"
	"testing"
	"time"
)

func mtime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func defaultFileMeta(inode string) Meta {
	return Meta{
		Owner: "alice", Group: "staff",
		Mode: 0o644, ModeStr: "-rw-r--r--",
		Inode: inode, Nlinks: 1,
	}
}

func defaultDirMeta(inode string) Meta {
	return Meta{
		Owner: "alice", Group: "staff",
		Mode: 0o755, ModeStr: "drwxr-xr-x",
		Inode: inode, Nlinks: 2,
	}
}

func TestStatAndLookup(t *testing.T) {
	spec := Dir("root", defaultDirMeta("1"),
		File("a.txt", 100, mtime("2026-01-01 00:00:00"), defaultFileMeta("2")),
	)
	f := New(spec)

	fi, err := f.Stat("/root/a.txt")
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if fi.Name() != "a.txt" || fi.Size() != 100 || fi.IsDir() {
		t.Errorf("unexpected fileinfo: %#v", fi)
	}

	if _, err := f.Stat("/root/missing"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestReadDirOrderAndMetadata(t *testing.T) {
	spec := Dir("root", defaultDirMeta("1"),
		File("b.txt", 10, mtime("2026-01-02 00:00:00"), defaultFileMeta("3")),
		File("a.txt", 20, mtime("2026-01-01 00:00:00"), defaultFileMeta("2")),
	)
	f := New(spec)

	entries, err := f.ReadDir("/root")
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	if entries[0].Name() != "a.txt" || entries[1].Name() != "b.txt" {
		t.Errorf("entries not name-sorted: %s, %s", entries[0].Name(), entries[1].Name())
	}
}

func TestUnreadableDirectory(t *testing.T) {
	spec := Dir("root", defaultDirMeta("1"),
		Unreadable(Dir("locked", defaultDirMeta("2"))),
	)
	f := New(spec)
	if _, err := f.ReadDir("/root/locked"); err == nil {
		t.Error("expected EACCES")
	}
}

func TestSymlinkResolution(t *testing.T) {
	spec := Dir("root", defaultDirMeta("1"),
		File("target.txt", 50, mtime("2026-01-01 00:00:00"), defaultFileMeta("2")),
		Symlink("link", "target.txt", defaultFileMeta("3")),
	)
	f := New(spec)

	fi, err := f.Lstat("/root/link")
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}
	if fi.Mode()&fs.ModeSymlink == 0 {
		t.Error("lstat should report symlink")
	}
	if !f.IsLink("/root/link") {
		t.Error("IsLink should return true")
	}

	resolved, err := f.EvalSymlinks("/root/link")
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	if resolved != "/root/target.txt" {
		t.Errorf("eval returned %q", resolved)
	}

	if got := f.Indicator("/root/link", false); got != "@" {
		t.Errorf("short indicator = %q", got)
	}
	if got := f.Indicator("/root/link", true); got != " ~> /root/target.txt" {
		t.Errorf("long indicator = %q", got)
	}
}

func TestOwnerGroupModeBlocks(t *testing.T) {
	spec := Dir("root", defaultDirMeta("1"),
		File("a.txt", 1024, mtime("2026-01-01 00:00:00"), defaultFileMeta("2")),
	)
	f := New(spec)
	fi, _ := f.Stat("/root/a.txt")

	owner, group := f.OwnerGroup(fi, true, true)
	if owner != "alice" || group != " staff  " {
		t.Errorf("owner/group = %q/%q", owner, group)
	}
	owner, group = f.OwnerGroup(fi, true, false)
	if owner != "alice" || group != "" {
		t.Errorf("owner-only: %q/%q", owner, group)
	}
	if got := f.ModeExtended(fi, "/root/a.txt"); got != "-rw-r--r-- " {
		t.Errorf("mode = %q", got)
	}
	if got := f.Blocks(fi); got != 2 {
		t.Errorf("blocks = %d", got)
	}
}

func TestPathOps(t *testing.T) {
	f := New(Dir("root", defaultDirMeta("1")))
	if got := f.Separator(); got != "/" {
		t.Errorf("separator = %q", got)
	}
	if got := f.Join("a", "b", "c"); got != "a/b/c" {
		t.Errorf("join = %q", got)
	}
	if got := f.Base("/root/a.txt"); got != "a.txt" {
		t.Errorf("base = %q", got)
	}
	if got := f.Ext("a.txt"); got != ".txt" {
		t.Errorf("ext = %q", got)
	}
	rel, _ := f.Rel("/root", "/root/sub/x")
	if rel != "sub/x" {
		t.Errorf("rel = %q", rel)
	}
}

func TestGitStatus(t *testing.T) {
	spec := Dir("root", defaultDirMeta("1"),
		File("a.txt", 1, mtime("2026-01-01 00:00:00"), defaultFileMeta("2")),
	)
	f := New(spec, WithGitStatus(map[string]string{"a.txt": "M"}))
	got := f.GitStatus("/root")
	if got["a.txt"] != "M" {
		t.Errorf("git status = %v", got)
	}
}
