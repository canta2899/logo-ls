package app

import (
	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/model"
)

type Config struct {
	FileList        []string
	AllMode         model.Include
	SortMode        model.SortMode
	LongListingMode model.Listing
	TimeFormatter   format.Timestamp
	Recursive       bool
	GitStatus       bool
	Reverse         bool
	DisableIcon     bool
	OneFilePerLine  bool
	Directory       bool
	NoGroup         bool
	HumanReadable   bool
	ShowBlockSize   bool
	ShowInodeNumber bool
}

func NewConfig() *Config {

	return &Config{
		AllMode:         model.IncludeDefault,
		SortMode:        model.SortAlphabetical,
		LongListingMode: model.LongListingNone,
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
