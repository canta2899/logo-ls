// Package cli parses command-line flags and produces the Config that the
// rest of logo-ls reads from. It owns nothing that touches the filesystem;
// app.App wires the parsed config to an inspector + renderer.
package cli

type Config struct {
	FileList          []string
	AllMode           Include
	SortMode          SortMode
	LongListingMode   Listing
	TimeFormatter     Timestamp
	Recursive         bool
	GitStatus         bool
	Reverse           bool
	DisableIcon       bool
	OneFilePerLine    bool
	Directory         bool
	NoGroup           bool
	HumanReadable     bool
	ShowBlockSize     bool
	ShowInodeNumber   bool
	NoIconOverride   bool
	IconOverrideFile string
}

func NewConfig() *Config {
	return &Config{
		AllMode:         IncludeDefault,
		SortMode:        SortAlphabetical,
		LongListingMode: LongListingNone,
		Recursive:       false,
		GitStatus:       false,
		Reverse:         false,
		DisableIcon:     false,
		OneFilePerLine:  false,
		Directory:       false,
		NoGroup:         false,
		HumanReadable:   false,
		ShowBlockSize:   false,
		ShowInodeNumber: false,
		TimeFormatter:   nil,
		FileList:        []string{},
	}
}
