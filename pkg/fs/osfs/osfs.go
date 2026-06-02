// Package osfs is the actual fs.FS implementation (which is OS dependent)
package osfs

import (
	iofs "io/fs"
	"os"
	"path/filepath"

	"github.com/canta2899/logo-ls/internal/inspect/git"
	"github.com/canta2899/logo-ls/pkg/fs"
)

func New() fs.FS {
	return &osFS{}
}

type osFS struct{}

func (o *osFS) Abs(path string) (string, error) { return filepath.Abs(path) }

func (o *osFS) Open(path string) (fs.File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &osFile{File: f}, nil
}

func (o *osFS) Stat(path string) (fs.FileInfo, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &osFileInfo{FileInfo: fi}, nil
}

func (o *osFS) Lstat(path string) (fs.FileInfo, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	return &osFileInfo{FileInfo: fi}, nil
}

func (o *osFS) ReadDir(path string) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil && len(entries) == 0 {
		return nil, err
	}
	out := make([]fs.DirEntry, 0, len(entries))
	for _, de := range entries {
		out = append(out, &osDirEntry{DirEntry: de})
	}
	return out, err
}

func (o *osFS) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (o *osFS) Join(parts ...string) string { return filepath.Join(parts...) }

func (o *osFS) Separator() string { return string(os.PathSeparator) }

func (o *osFS) Base(path string) string { return filepath.Base(path) }

func (o *osFS) Dir(path string) string { return filepath.Dir(path) }

func (o *osFS) Ext(path string) string { return filepath.Ext(path) }

func (o *osFS) Rel(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

func (o *osFS) FromSlash(path string) string {
	return filepath.FromSlash(path)
}

var osGitReader = git.NewStatusReader(git.ExecPorcelain{})

func (o *osFS) GitStatus(dir string) map[string]string {
	m := osGitReader.StatusRelative(dir)
	if len(m) == 0 {
		return nil
	}
	return m
}

// osFile wraps *os.File to satisfy fs.File.
type osFile struct {
	*os.File
}

func (f *osFile) Stat() (fs.FileInfo, error) {
	fi, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return &osFileInfo{FileInfo: fi}, nil
}

func (f *osFile) ReadDir(n int) ([]fs.DirEntry, error) {
	entries, err := f.File.ReadDir(n)
	out := make([]fs.DirEntry, 0, len(entries))
	for _, de := range entries {
		out = append(out, &osDirEntry{DirEntry: de})
	}
	return out, err
}

// osFileInfo wraps os.FileInfo. It exists so we can swap implementations
// without exposing concrete os.FileInfo types throughout the codebase.
type osFileInfo struct {
	os.FileInfo
}

// osDirEntry wraps os.DirEntry and converts Info() to our type.
type osDirEntry struct {
	iofs.DirEntry
}

func (d *osDirEntry) Info() (fs.FileInfo, error) {
	fi, err := d.DirEntry.Info()
	if fi == nil {
		return nil, err
	}
	return &osFileInfo{FileInfo: fi}, err
}
