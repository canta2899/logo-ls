package app

import (
	"time"

	"github.com/canta2899/logo-ls/internal/model"
)

type Config struct {
	FileList        []model.FileEntry
	AllMode         model.Include
	SortMode        model.SortMode
	LongListingMode model.Listing
	TimeFormat      string
	Recursive       bool
	GitStatus       bool
	Reverse         bool
	DisableColor    bool
	DisableIcon     bool
	OneFilePerLine  bool
	Directory       bool
	NoGroup         bool
	HumanReadable   bool
	ShowBlockSize   bool
}

func NewConfig() *Config {

	return &Config{
		AllMode:         model.IncludeDefault,
		SortMode:        model.SortAlphabetical,
		LongListingMode: model.LongListingNone,
		Recursive:       false,
		GitStatus:       false,
		Reverse:         false,
		DisableColor:    false,
		DisableIcon:     false,
		OneFilePerLine:  false,
		Directory:       false,
		NoGroup:         false,
		HumanReadable:   false,
		ShowBlockSize:   false,
		TimeFormat:      time.Stamp,
		FileList:        []model.FileEntry{},
	}
}
