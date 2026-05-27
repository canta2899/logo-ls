// Package fs defines a filesystem abstraction used by the rest of logo-ls.
package fs

import (
	iofs "io/fs"
)

type FS interface {
	Abs(path string) (string, error)
	Open(path string) (File, error)
	Stat(path string) (FileInfo, error)
	Lstat(path string) (FileInfo, error)
	ReadDir(path string) ([]DirEntry, error)
	EvalSymlinks(path string) (string, error)

	Join(parts ...string) string
	Separator() string
	Base(path string) string
	Dir(path string) string
	Ext(path string) string
	Rel(base, target string) (string, error)
	FromSlash(path string) string

	// Indicator returns the suffix character appended to a name ("/", "@",
	Indicator(path string, longMode bool) string

	IsLink(path string) bool

	InodeNumber(path string) string

	HardLinks(path string) uint64

	ModeExtended(fi FileInfo, path string) string

	OwnerGroup(fi FileInfo, showOwner, showGroup bool) (owner, group string)

	Blocks(fi FileInfo) int64

	GitStatus(dir string) map[string]string
}

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
