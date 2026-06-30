package git

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractStatusChar(t *testing.T) {
	cases := []struct {
		xy, want string
	}{
		{"??", "U"},
		{"!!", "I"},
		{"M ", "M"},
		{" M", "M"},
		{"A ", "A"},
		{"R", "R"},
		{"UU", "U"},
		{"  ", "?"},
		{"", "?"},
	}
	for _, c := range cases {
		t.Run(c.xy, func(t *testing.T) {
			if got := ExtractStatusChar(c.xy); got != c.want {
				t.Errorf("ExtractStatusChar(%q) = %q, want %q", c.xy, got, c.want)
			}
		})
	}
}

func TestParsePorcelain_Empty(t *testing.T) {
	m := ParsePorcelain("/repo", nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

// helper: build a -z porcelain stream from line strings.
func zstream(lines ...string) []byte {
	return []byte(strings.Join(lines, "\000") + "\000")
}

func TestParsePorcelain_Modified(t *testing.T) {
	raw := zstream(" M src/file.go")
	got := ParsePorcelain("/repo", raw)

	wantFile := filepath.Clean("/repo/src/file.go")
	if got[wantFile] != "M" {
		t.Errorf("missing or wrong status for %s: %v", wantFile, got)
	}
	// Parent directory should be marked with "M"
	wantDir := filepath.Clean("/repo/src") + string(filepath.Separator)
	if got[wantDir] != "M" {
		t.Errorf("parent dir not marked: %v", got)
	}
}

func TestParsePorcelain_Untracked(t *testing.T) {
	raw := zstream("?? newfile.txt")
	got := ParsePorcelain("/repo", raw)
	want := filepath.Clean("/repo/newfile.txt")
	if got[want] != "U" {
		t.Errorf("untracked = %q, want U", got[want])
	}
}

func TestParsePorcelain_Staged(t *testing.T) {
	raw := zstream("A  added.txt", "M  edited.txt")
	got := ParsePorcelain("/repo", raw)
	if got[filepath.Clean("/repo/added.txt")] != "A" {
		t.Errorf("added: %v", got)
	}
	if got[filepath.Clean("/repo/edited.txt")] != "M" {
		t.Errorf("edited: %v", got)
	}
}

func TestParsePorcelain_Unmerged(t *testing.T) {
	raw := zstream("UU conflict.txt")
	got := ParsePorcelain("/repo", raw)
	if got[filepath.Clean("/repo/conflict.txt")] != "U" {
		t.Errorf("unmerged: %v", got)
	}
}

func TestParsePorcelain_NestedParents(t *testing.T) {
	raw := zstream(" M a/b/c/file.go")
	got := ParsePorcelain("/repo", raw)

	for _, p := range []string{"/repo/a", "/repo/a/b", "/repo/a/b/c"} {
		key := filepath.Clean(p) + string(filepath.Separator)
		if got[key] != "M" {
			t.Errorf("parent %s not marked: %v", key, got)
		}
	}
}

func TestParsePorcelain_DoesNotOverwriteFile(t *testing.T) {
	// If a file's parent is also explicitly modified by another file beneath
	// it, parent stays "M" and the explicit file status is preserved.
	raw := zstream(" M dir/a.go", " M dir")
	got := ParsePorcelain("/repo", raw)

	// The explicit "dir" entry from porcelain wins over the auto-parent mark
	// at file level. Auto-parent only fills slots that don't already exist.
	if got[filepath.Clean("/repo/dir/a.go")] != "M" {
		t.Errorf("file overwritten: %v", got)
	}
}

func TestParsePorcelain_Ignored(t *testing.T) {
	raw := zstream("!! ignored.log")
	got := ParsePorcelain("/repo", raw)
	want := filepath.Clean("/repo/ignored.log")
	if got[want] != "I" {
		t.Errorf("ignored = %q, want I", got[want])
	}
}

func TestParsePorcelain_Rename(t *testing.T) {
	// Porcelain -z renames: `R<space>newpath\0oldpath`.
	raw := []byte("R  new/path.go\x00old/path.go\x00")
	got := ParsePorcelain("/repo", raw)
	if got[filepath.Clean("/repo/new/path.go")] != "R" {
		t.Errorf("rename new path missing: %v", got)
	}
	// The old path should NOT be present as a separate record.
	if _, ok := got[filepath.Clean("/repo/old/path.go")]; ok {
		t.Errorf("rename old path leaked as record: %v", got)
	}
}

func TestParsePorcelain_PathWithLeadingSpace(t *testing.T) {
	// XY is exactly two chars then a single space then path - a path that
	// itself starts with a space must be preserved.
	raw := zstream("A   spaced.txt")
	got := ParsePorcelain("/repo", raw)
	want := filepath.Clean("/repo/ spaced.txt")
	if got[want] != "A" {
		t.Errorf("expected leading-space path preserved, got %v", got)
	}
}

// fakePorcelain implements the Porcelain interface for testing the
// StatusReader caching behaviour.
type fakePorcelain struct {
	root    string
	raw     []byte
	calls   int
	rootErr error
}

func (f *fakePorcelain) Root(dir string) (string, error) {
	if f.rootErr != nil {
		return "", f.rootErr
	}
	return f.root, nil
}

func (f *fakePorcelain) Status(root string) ([]byte, error) {
	f.calls++
	return f.raw, nil
}

func TestStatusReader_Caches(t *testing.T) {
	fp := &fakePorcelain{
		root: "/repo",
		raw:  zstream(" M file.go"),
	}
	r := NewStatusReader(fp)

	_ = r.Status("/repo/sub")
	_ = r.Status("/repo/other")

	if fp.calls != 1 {
		t.Errorf("expected porcelain Status to be called once, got %d", fp.calls)
	}
}

func TestStatusReader_NotARepo(t *testing.T) {
	fp := &fakePorcelain{rootErr: &noRepoErr{}}
	r := NewStatusReader(fp)
	if m := r.Status("/somewhere"); m != nil {
		t.Errorf("expected nil for non-repo, got %v", m)
	}
}

type noRepoErr struct{}

func (noRepoErr) Error() string { return "not a git repository" }
