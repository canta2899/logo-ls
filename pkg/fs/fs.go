// Package fs defines a filesystem abstraction used by the rest of logo-ls.
package fs

import (
	iofs "io/fs"
)

// FS is the narrow filesystem abstraction. It exposes only real OS
// operations and a small set of path manipulation helpers that vary by
// platform (separator, FromSlash). Higher-level concerns — indicator
// strings, owner/group resolution, git status, mode formatting — live in
// internal/inspect and internal/render.
type FS interface {
	// Real filesystem operations.
	Abs(path string) (string, error)
	Open(path string) (File, error)
	Stat(path string) (FileInfo, error)
	Lstat(path string) (FileInfo, error)
	ReadDir(path string) ([]DirEntry, error)
	EvalSymlinks(path string) (string, error)

	// Path manipulation. Separator and FromSlash are OS-dependent.
	Join(parts ...string) string
	Separator() string
	Base(path string) string
	Dir(path string) string
	Ext(path string) string
	Rel(base, target string) (string, error)
	FromSlash(path string) string

	// GitStatus is a legacy convenience for callers that haven't yet been
	// migrated to git.StatusReader. New code should use the reader directly.
	// Returns nil when path is not inside a git repository.
	GitStatus(dir string) map[string]string
}

// FileMode aliases io/fs.FileMode for callers in this module's packages.
type FileMode = iofs.FileMode

type FileInfo interface {
	iofs.FileInfo
}

type DirEntry interface {
	Name() string
	IsDir() bool
	Type() iofs.FileMode
	Info() (FileInfo, error)
}

type File interface {
	Name() string
	Stat() (FileInfo, error)
	ReadDir(n int) ([]DirEntry, error)
	Close() error
}
