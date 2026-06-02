package inspect_test

import (
	"testing"
	"time"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
	"github.com/canta2899/logo-ls/internal/inspect"
)

func mtime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestInspector_Basic(t *testing.T) {
	root := fakefs.Dir("root", fakefs.Meta{
		Owner: "alice", Group: "staff", Mode: 0o755,
		Inode: "100", Nlinks: 2, Mtime: mtime("2026-01-01 00:00:00"),
	},
		fakefs.File("readme.md", 42, mtime("2026-01-02 10:00:00"),
			fakefs.Meta{
				Owner: "alice", Group: "staff", Mode: 0o644,
				Inode: "200", Nlinks: 1, Blocks: 8,
			}),
	)
	vfs := fakefs.New(root)

	fi, err := vfs.Lstat("/root/readme.md")
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}

	insp := inspect.New(vfs, inspect.DefaultIconResolver(), inspect.Options{
		Long:       true,
		ShowOwner:  true,
		ShowGroup:  true,
		ShowInode:  true,
		ShowBlocks: true,
	})
	e := insp.Inspect("/root/readme.md", fi)

	if e.Name != "readme.md" {
		t.Errorf("name: %s", e.Name)
	}
	if e.Owner != "alice" {
		t.Errorf("owner: %q", e.Owner)
	}
	if e.Group != "staff" {
		t.Errorf("group: %q", e.Group)
	}
	if e.Inode != "200" {
		t.Errorf("inode: %q", e.Inode)
	}
	if e.HardLinks != 1 {
		t.Errorf("hardlinks: %d", e.HardLinks)
	}
	if e.Blocks != 8 {
		t.Errorf("blocks: %d", e.Blocks)
	}
	if e.Kind != inspect.KindFile {
		t.Errorf("kind: %v", e.Kind)
	}
}

func TestInspector_NoLongModeSkipsOwnerLookup(t *testing.T) {
	root := fakefs.Dir("root", fakefs.Meta{Mode: 0o755},
		fakefs.File("foo", 1, mtime("2026-01-02 10:00:00"),
			fakefs.Meta{Owner: "alice", Group: "staff", Mode: 0o644}),
	)
	vfs := fakefs.New(root)

	fi, err := vfs.Lstat("/root/foo")
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}
	insp := inspect.New(vfs, inspect.DefaultIconResolver(), inspect.Options{Long: false})
	e := insp.Inspect("/root/foo", fi)

	if e.Owner != "" || e.Group != "" {
		t.Errorf("non-long mode should not populate owner/group, got %q/%q", e.Owner, e.Group)
	}
	if e.HardLinks != 0 {
		t.Errorf("non-long mode should not populate hardlinks, got %d", e.HardLinks)
	}
}

func TestInspector_DirIndicator(t *testing.T) {
	root := fakefs.Dir("root", fakefs.Meta{Mode: 0o755},
		fakefs.Dir("sub", fakefs.Meta{Mode: 0o755}),
	)
	vfs := fakefs.New(root)
	fi, err := vfs.Lstat("/root/sub")
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}
	insp := inspect.New(vfs, inspect.DefaultIconResolver(), inspect.Options{})
	e := insp.Inspect("/root/sub", fi)
	if e.Indicator != "/" {
		t.Errorf("dir indicator: %q", e.Indicator)
	}
	if e.Kind != inspect.KindDir {
		t.Errorf("kind: %v", e.Kind)
	}
}
