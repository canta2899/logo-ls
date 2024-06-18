package model

type Include int

const (
	IncludeDefault Include = iota
	IncludeAlmost
	IncludeAll
)

type SortMode int

const (
	SortSize SortMode = iota
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
