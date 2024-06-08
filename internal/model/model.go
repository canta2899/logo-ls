package model

import "path/filepath"

type Include int

const (
	IncludeDefault Include = iota
	IncludeAlmost
	IncludeAll
)

type Sort int

const (
	SortSize Sort = iota
	SortModTime
	SortExtension
	SortAlphabetical
	SortNatural
	SortNone
)

type Listing int

const (
	LongListingOwner Listing = iota
	LongListingGroup
	LongListingDefault
	LongListingNone
)

type ExitCode int

const (
	CodeOk ExitCode = iota
	CodeMinor
	CodeSerious
)

func (e *ExitCode) SetMinor() {
	if *e == CodeSerious {
		return
	}

	*e = CodeMinor
}

func (e *ExitCode) SetSerious() {
	*e = CodeSerious
}

type FileEntry struct {
	Path         string
	AbsolutePath string
}

func NewFileEntry(path string) FileEntry {

	abs, err := filepath.Abs(path)

	if err != nil {
		abs = path
	}

	return FileEntry{
		Path:         path,
		AbsolutePath: abs,
	}
}
