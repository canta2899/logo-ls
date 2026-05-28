// Package fakefs is an in-memory fs.FS implementation for tests.
package fakefs

import (
	"errors"
	"fmt"
	iofs "io/fs"
	"maps"
	"path"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/canta2899/logo-ls/fs"
)

type Meta struct {
	Owner   string
	Group   string
	Mode    iofs.FileMode
	ModeStr string
	Inode   string
	Nlinks  uint64
	Blocks  int64
	Mtime   time.Time
}

type EntryKind int

const (
	kindFile EntryKind = iota
	kindDir
	kindSymlink
)

type Entry struct {
	kind     EntryKind
	name     string
	size     int64
	mtime    time.Time
	meta     Meta
	target   string   // for symlinks
	children []*Entry // for dirs

	// Failures: makes ReadDir return EACCES on this directory.
	unreadable bool
}

// File creates a file entry.
func File(name string, size int64, mtime time.Time, meta Meta) *Entry {
	return &Entry{
		kind:  kindFile,
		name:  name,
		size:  size,
		mtime: mtime,
		meta:  meta,
	}
}

// Dir creates a directory entry.
func Dir(name string, meta Meta, children ...*Entry) *Entry {
	if meta.Mode == 0 {
		meta.Mode = iofs.ModeDir | 0o755
	} else {
		meta.Mode |= iofs.ModeDir
	}
	return &Entry{
		kind:     kindDir,
		name:     name,
		mtime:    meta.Mtime,
		meta:     meta,
		children: children,
	}
}

// Symlink creates a symlink entry pointing at target.
func Symlink(name, target string, meta Meta) *Entry {
	meta.Mode |= iofs.ModeSymlink
	return &Entry{
		kind:   kindSymlink,
		name:   name,
		target: target,
		mtime:  meta.Mtime,
		meta:   meta,
	}
}

// Unreadable marks a directory as returning EACCES on ReadDir.
func Unreadable(d *Entry) *Entry {
	d.unreadable = true
	return d
}

// Option configures the fake filesystem.
type Option func(*fakeFS)

// WithGitStatus sets the map returned by GitStatus(dir).
func WithGitStatus(status map[string]string) Option {
	return func(f *fakeFS) { f.gitStatus = status }
}

// New returns a fs.FS backed by the given spec rooted at "/".
//
// The provided root entry is mounted at "/<root.name>" (so paths look like
// "/root/foo.txt"). Use the same root name as your test's input path.
func New(root *Entry, opts ...Option) fs.FS {
	f := &fakeFS{
		root: root,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

type fakeFS struct {
	root      *Entry
	gitStatus map[string]string
}

func (f *fakeFS) Join(parts ...string) string { return path.Join(parts...) }

func (f *fakeFS) Separator() string { return "/" }

func (f *fakeFS) Base(p string) string { return path.Base(p) }

func (f *fakeFS) Dir(p string) string { return path.Dir(p) }

func (f *fakeFS) Ext(p string) string { return path.Ext(p) }

func (f *fakeFS) FromSlash(p string) string { return p }

func (f *fakeFS) Rel(base, target string) (string, error) {
	base = path.Clean(base)
	target = path.Clean(target)
	if base == target {
		return ".", nil
	}
	if strings.HasPrefix(target, base+"/") {
		return target[len(base)+1:], nil
	}
	return target, nil
}

func (f *fakeFS) Abs(p string) (string, error) {
	if strings.HasPrefix(p, "/") {
		return path.Clean(p), nil
	}
	return path.Clean("/" + p), nil
}

func (f *fakeFS) lookup(p string) (*Entry, error) {
	p = path.Clean(p)
	if p == "/" {
		// Synthetic root above the mounted entry. Returns dir metadata so
		// listing ".." from /root renders cleanly.
		rootMtime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		return &Entry{
			kind:  kindDir,
			name:  "/",
			mtime: rootMtime,
			meta: Meta{
				Owner: "root", Group: "wheel",
				Mode: iofs.ModeDir | 0o755, ModeStr: "drwxr-xr-x",
				Inode: "1", Nlinks: 1, Mtime: rootMtime,
			},
			children: []*Entry{f.root},
		}, nil
	}
	parts := strings.Split(strings.TrimPrefix(p, "/"), "/")
	if parts[0] != f.root.name {
		return nil, &iofs.PathError{Op: "stat", Path: p, Err: iofs.ErrNotExist}
	}
	cur := f.root
	for _, part := range parts[1:] {
		if cur.kind != kindDir {
			return nil, &iofs.PathError{Op: "stat", Path: p, Err: iofs.ErrNotExist}
		}
		var next *Entry
		for _, c := range cur.children {
			if c.name == part {
				next = c
				break
			}
		}
		if next == nil {
			return nil, &iofs.PathError{Op: "stat", Path: p, Err: iofs.ErrNotExist}
		}
		cur = next
	}
	return cur, nil
}

// resolveSymlink follows a symlink relative to its parent directory.
func (f *fakeFS) resolveSymlink(linkPath string, e *Entry) (string, *Entry, error) {
	target := e.target
	if !path.IsAbs(target) {
		target = path.Join(path.Dir(linkPath), target)
	}
	target = path.Clean(target)
	resolved, err := f.lookup(target)
	if err != nil {
		return target, nil, err
	}
	if resolved.kind == kindSymlink {
		return f.resolveSymlink(target, resolved)
	}
	return target, resolved, nil
}

func (f *fakeFS) Stat(p string) (fs.FileInfo, error) {
	e, err := f.lookup(p)
	if err != nil {
		return nil, err
	}
	if e.kind == kindSymlink {
		_, target, err := f.resolveSymlink(p, e)
		if err != nil {
			return nil, err
		}
		return &fakeFileInfo{entry: target}, nil
	}
	return &fakeFileInfo{entry: e}, nil
}

func (f *fakeFS) Lstat(p string) (fs.FileInfo, error) {
	e, err := f.lookup(p)
	if err != nil {
		return nil, err
	}
	return &fakeFileInfo{entry: e}, nil
}

func (f *fakeFS) Open(p string) (fs.File, error) {
	e, err := f.lookup(p)
	if err != nil {
		return nil, err
	}
	if e.kind == kindSymlink {
		_, target, err := f.resolveSymlink(p, e)
		if err != nil {
			return nil, err
		}
		e = target
	}
	return &fakeFile{fs: f, entry: e, absPath: path.Clean(p)}, nil
}

func (f *fakeFS) ReadDir(p string) ([]fs.DirEntry, error) {
	e, err := f.lookup(p)
	if err != nil {
		return nil, err
	}
	if e.kind != kindDir {
		return nil, &iofs.PathError{Op: "readdir", Path: p, Err: errors.New("not a directory")}
	}
	if e.unreadable {
		return nil, &iofs.PathError{Op: "readdir", Path: p, Err: syscall.EACCES}
	}
	// Return entries in name order so test output is stable.
	children := append([]*Entry(nil), e.children...)
	sort.Slice(children, func(i, j int) bool { return children[i].name < children[j].name })
	out := make([]fs.DirEntry, 0, len(children))
	for _, c := range children {
		out = append(out, &fakeDirEntry{entry: c})
	}
	return out, nil
}

func (f *fakeFS) EvalSymlinks(p string) (string, error) {
	e, err := f.lookup(p)
	if err != nil {
		return "", err
	}
	if e.kind != kindSymlink {
		return path.Clean(p), nil
	}
	target, _, err := f.resolveSymlink(p, e)
	if err != nil {
		return "", err
	}
	return target, nil
}

func (f *fakeFS) Indicator(p string, longMode bool) string {
	e, err := f.lookup(p)
	if err != nil {
		return ""
	}
	mode := e.meta.Mode
	switch {
	case mode&iofs.ModeDir > 0:
		return "/"
	case mode&iofs.ModeNamedPipe > 0:
		return "|"
	case mode&iofs.ModeSymlink > 0:
		if !longMode {
			return "@"
		}
		target, _, err := f.resolveSymlink(p, e)
		if err != nil {
			return ""
		}
		return " ~> " + target
	case mode&iofs.ModeSocket > 0:
		return "="
	case mode&0o111 != 0:
		return "*"
	}
	return ""
}

func (f *fakeFS) IsLink(p string) bool {
	e, err := f.lookup(p)
	if err != nil {
		return false
	}
	return e.meta.Mode&iofs.ModeSymlink > 0
}

func (f *fakeFS) InodeNumber(p string) string {
	e, err := f.lookup(p)
	if err != nil {
		return ""
	}
	return e.meta.Inode
}

func (f *fakeFS) HardLinks(p string) uint64 {
	e, err := f.lookup(p)
	if err != nil {
		return 0
	}
	if e.meta.Nlinks == 0 {
		return 1
	}
	return e.meta.Nlinks
}

func (f *fakeFS) ModeExtended(fi fs.FileInfo, p string) string {
	if fi == nil {
		return fmt.Sprintf("%-*s", 11, "???????????")
	}
	if ff, ok := fi.(*fakeFileInfo); ok && ff.entry.meta.ModeStr != "" {
		return fmt.Sprintf("%-*s", 11, ff.entry.meta.ModeStr)
	}
	return fmt.Sprintf("%-*s", 11, fi.Mode().String())
}

func (f *fakeFS) OwnerGroup(fi fs.FileInfo, showOwner, showGroup bool) (string, string) {
	if fi == nil {
		return "", ""
	}
	ff, ok := fi.(*fakeFileInfo)
	if !ok {
		return "", ""
	}
	var owner, group string
	if showOwner {
		owner = ff.entry.meta.Owner
	}
	if showGroup && ff.entry.meta.Group != "" {
		group = fmt.Sprintf(" %v  ", ff.entry.meta.Group)
	}
	return owner, group
}

func (f *fakeFS) Blocks(fi fs.FileInfo) int64 {
	if fi == nil {
		return 0
	}
	if ff, ok := fi.(*fakeFileInfo); ok {
		if ff.entry.meta.Blocks != 0 {
			return ff.entry.meta.Blocks
		}
		// Default: ceil(size/512).
		return (ff.entry.size + 511) / 512
	}
	return 0
}

func (f *fakeFS) GitStatus(dir string) map[string]string {
	if len(f.gitStatus) == 0 {
		return nil
	}
	out := make(map[string]string, len(f.gitStatus))
	maps.Copy(out, f.gitStatus)
	return out
}

type fakeFileInfo struct {
	entry *Entry
}

func (fi *fakeFileInfo) Name() string        { return fi.entry.name }
func (fi *fakeFileInfo) Size() int64         { return fi.entry.size }
func (fi *fakeFileInfo) Mode() iofs.FileMode { return fi.entry.meta.Mode }
func (fi *fakeFileInfo) ModTime() time.Time  { return fi.entry.mtime }
func (fi *fakeFileInfo) IsDir() bool         { return fi.entry.kind == kindDir }
func (fi *fakeFileInfo) Sys() any            { return nil }

type fakeDirEntry struct {
	entry *Entry
}

func (de *fakeDirEntry) Name() string { return de.entry.name }
func (de *fakeDirEntry) IsDir() bool  { return de.entry.kind == kindDir }
func (de *fakeDirEntry) Type() iofs.FileMode {
	return de.entry.meta.Mode.Type()
}

func (de *fakeDirEntry) Info() (fs.FileInfo, error) {
	return &fakeFileInfo{entry: de.entry}, nil
}

type fakeFile struct {
	fs      *fakeFS
	entry   *Entry
	absPath string
}

func (fl *fakeFile) Name() string { return fl.absPath }

func (fl *fakeFile) Stat() (fs.FileInfo, error) {
	return &fakeFileInfo{entry: fl.entry}, nil
}

func (fl *fakeFile) ReadDir(n int) ([]fs.DirEntry, error) {
	return fl.fs.ReadDir(fl.absPath)
}

func (fl *fakeFile) Close() error { return nil }
