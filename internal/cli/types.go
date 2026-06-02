package cli

// Include selects which directory entries to list.
type Include int

const (
	IncludeDefault Include = iota
	IncludeAlmost
	IncludeAll
)

// SortMode selects the comparator used to order entries.
type SortMode int

const (
	SortSize SortMode = iota
	SortModTime
	SortExtension
	SortAlphabetical
	SortNatural
	SortNone
)

// Listing selects the long-mode column set.
type Listing int

const (
	LongListingOwner Listing = iota
	LongListingGroup
	LongListingDefault
	LongListingNone
)

// ExitCode is the process exit status the CLI reports.
type ExitCode int

const (
	CodeOk ExitCode = iota
	CodeMinor
	CodeSerious
)

// SetMinor sets the exit code to CodeMinor unless a more severe code is
// already in effect.
func (e *ExitCode) SetMinor() {
	if *e == CodeSerious {
		return
	}
	*e = CodeMinor
}

// SetSerious unconditionally raises the exit code to CodeSerious.
func (e *ExitCode) SetSerious() { *e = CodeSerious }
