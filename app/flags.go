package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/model"
)

var Version = ""

// ErrHelp and ErrVersionRequested are returned by BuildConfig when the user
// requested --help or --version.
var (
	ErrHelp             = errors.New("help requested")
	ErrVersionRequested = errors.New("version requested")
)

// BuildConfig parses the given arg list (including argv[0]) and returns a
// Config and the parser used to build it. It never calls os.Exit, making it
// safe to use from tests.
//
// On --help or --version, it returns the sentinel errors above so the caller
// can choose how to render them. The parser is returned in all cases.
func BuildConfig(argv []string) (*Config, *Parser, error) {
	c := NewConfig()
	opt := NewParser()
	opt.Parameters = "[files ...]"

	help := opt.Bool('?', "help", "display this help and exit")
	version := opt.Bool('V', "version", "output version information and exit")

	includeAll := opt.Bool('a', "all", "do not ignore entries starting with .")
	includeAlmost := opt.Bool('A', "almost-all", "do not list implied . and ..")

	c.SortMode = model.SortAlphabetical

	sortNone := opt.Bool('U', "", "do not sort; list entries in directory order")
	sortNatural := opt.Bool('v', "", "natural sort of (version) numbers within text")
	sortExtension := opt.Bool('X', "", "sort alphabetically by entry extension")
	sortModTime := opt.Bool('t', "", "sort by modification time, newest first")
	sortSize := opt.Bool('S', "", "sort by file size, largest first")

	reverse := opt.Bool('r', "reverse", "reverse order while sorting")
	recursive := opt.Bool('R', "recursive", "list subdirectories recursively")
	gitStatus := opt.Bool('D', "git-status", "print git status of files")
	disableIcon := opt.Bool('e', "disable-icon", "don't print icons of the files")
	showInodeNumber := opt.Bool('i', "inode", "print the index number of each file")
	oneFilePerLine := opt.Bool('1', "", "list one file per line.")
	directory := opt.Bool('d', "directory", "list directories themselves, not their contents")
	noGroup := opt.Bool('G', "no-group", "in a long listing, don't print group names")
	humanReadable := opt.Bool('h', "human-readable", "with -l and -s, print sizes like 1K 234M 2G etc.")
	showBlockSize := opt.Bool('s', "size", "print the allocated size of each file, in blocks")

	completeTimeInformation := opt.Bool('T', "time-style", "display complete time information")

	c.LongListingMode = model.LongListingNone

	longListingMode := opt.Bool('o', "", "like -l, but do not list group information")
	longListingGroup := opt.Bool('g', "", "like -l, but do not list owner")
	longListingDefault := opt.Bool('l', "", "use a long listing format")

	if err := opt.Parse(argv); err != nil {
		return nil, opt, err
	}

	if *help {
		return nil, opt, ErrHelp
	}
	if *version {
		return nil, opt, ErrVersionRequested
	}

	switch {
	case *includeAll:
		c.AllMode = model.IncludeAll
	case *includeAlmost:
		c.AllMode = model.IncludeAlmost
	}

	switch {
	case *sortNone:
		c.SortMode = model.SortNone
	case *sortNatural:
		c.SortMode = model.SortNatural
	case *sortExtension:
		c.SortMode = model.SortExtension
	case *sortModTime:
		c.SortMode = model.SortModTime
	case *sortSize:
		c.SortMode = model.SortSize
	}

	switch {
	case *longListingMode:
		c.LongListingMode = model.LongListingOwner
	case *longListingGroup:
		c.LongListingMode = model.LongListingGroup
	case *longListingDefault:
		c.LongListingMode = model.LongListingDefault
	}

	c.TimeFormatter = format.GetFormatter(*completeTimeInformation)

	c.Reverse = *reverse
	c.Recursive = *recursive
	c.GitStatus = *gitStatus
	c.DisableIcon = *disableIcon
	c.OneFilePerLine = *oneFilePerLine
	c.Directory = *directory
	c.NoGroup = *noGroup
	c.HumanReadable = *humanReadable
	c.ShowBlockSize = *showBlockSize
	c.ShowInodeNumber = *showInodeNumber

	if len(opt.Args) > 0 {
		c.FileList = append(c.FileList, opt.Args...)
	} else {
		c.FileList = append(c.FileList, ".")
	}

	return c, opt, nil
}

// GetConfigFromCli parses os.Args. On --help/--version it prints to stdout
// and exits 0; on parse errors it prints to stderr and exits 2.
func GetConfigFromCli() *Config {
	c, opt, err := BuildConfig(os.Args)
	if err != nil {
		switch {
		case errors.Is(err, ErrHelp):
			printHelpMessage(opt)
			os.Exit(0)
		case errors.Is(err, ErrVersionRequested):
			printVersion()
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "logo-ls: %v\nTry 'logo-ls --help' for more information.\n", err)
			os.Exit(2)
		}
	}
	return c
}

func printVersion() {
	fmt.Printf("logo-ls %s\nLicense MIT <https://opensource.org/licenses/MIT>.\nThis is free software: you are free to change and redistribute it.\nThere is NO WARRANTY, to the extent permitted by law.\n", Version)
}

func printHelpMessage(opt *Parser) {
	fmt.Println("logo-ls: A modern ls command with icons and Git status integration.")
	fmt.Println("Lists information about the FILEs (the current directory by default).")
	fmt.Println("Sorts entries alphabetically if none of -tvSUX is specified.")
	fmt.Println()
	if opt != nil {
		opt.PrintUsage()
	}
}
