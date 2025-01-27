package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/canta2899/logo-ls/app"
	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/model"
	"github.com/pborman/getopt/v2"
	"golang.org/x/term"
)

const standardTerminalWidth = 80

func GetConfig() *app.Config {

	c := app.NewConfig()

	getopt.AllowAnyOrder(true)
	getopt.SetParameters("[files ...]")

	help := getopt.BoolLong("help", '?', "display this help and exit")
	version := getopt.BoolLong("version", 'V', "output version information and exit")

	includeAll := getopt.BoolLong("all", 'a', "do not ignore entries starting with .")
	includeAlmost := getopt.BoolLong("almost-all", 'A', "do not list implied . and ..")

	c.SortMode = model.SortAlphabetical

	sortNone := getopt.Bool('U', "do not sort; list entries in directory order")
	sortNatural := getopt.Bool('v', "natural sort of (version) numbers within text")
	sortExtension := getopt.Bool('X', "sort alphabetically by entry extension")
	sortModTime := getopt.Bool('t', "sort by modification time, newest first")
	sortSize := getopt.Bool('S', "sort by file size, largest first")

	reverse := getopt.BoolLong("reverse", 'r', "reverse order while sorting")
	recursive := getopt.BoolLong("recursive", 'R', "list subdirectories recursively")
	gitStatus := getopt.BoolLong("git-status", 'D', "print git status of files")
	disableColor := getopt.BoolLong("disable-color", 'c', "don't color icons, filenames and git status (use this to print to a file)")
	disableIcon := getopt.BoolLong("disable-icon", 'e', "don't print icons of the files")
	showInodeNumber := getopt.BoolLong("inode", 'i', "print the index number of each file")
	oneFilePerLine := getopt.Bool('1', "list one file per line.")
	directory := getopt.BoolLong("directory", 'd', "list directories themselves, not their contents")
	noGroup := getopt.BoolLong("no-group", 'G', "in a long listing, don't print group names")
	humanReadable := getopt.BoolLong("human-readable", 'h', "with -l and -s, print sizes like 1K 234M 2G etc.")
	showBlockSize := getopt.BoolLong("size", 's', "print the allocated size of each file, in blocks")

	completeTimeInformation := getopt.BoolLong("time-style", 'T', "display complete time information")

	c.LongListingMode = model.LongListingNone

	longListingMode := getopt.Bool('o', "like -l, but do not list group information")
	longListingGroup := getopt.Bool('g', "\nlike -l, but do not list owner")
	longListingDefault := getopt.Bool('l', "use a long listing format")

	// using getopt.Getopt instead of parse to provide custom err
	err := getopt.Getopt(nil)
	if err != nil {
		fmt.Printf("%v\nTry 'logo-ls -?' for more information.", err)
		os.Exit(1)
	}

	if *help {
		printHelpMessage()
		os.Exit(0)
	}

	if *version {
		printVersion()
		os.Exit(0)
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
	c.DisableColor = *disableColor
	c.DisableIcon = *disableIcon
	c.OneFilePerLine = *oneFilePerLine
	c.Directory = *directory
	c.NoGroup = *noGroup
	c.HumanReadable = *humanReadable
	c.ShowBlockSize = *showBlockSize
	c.ShowInodeNumber = *showInodeNumber

	args := getopt.Args()
	if len(args) > 0 {
		c.FileList = append(c.FileList, args...)
	} else {
		c.FileList = append(c.FileList, ".")
	}

	if c.LongListingMode == model.LongListingNone {
		c.TerminalWidth = getCustomTerminalWidth()
	} else {
		c.TerminalWidth = standardTerminalWidth
	}

	return c
}

func printVersion() {
	fmt.Printf("logo-ls %s\nLicense MIT <https://opensource.org/licenses/MIT>.\nThis is free software: you are free to change and redistribute it.\nThere is NO WARRANTY, to the extent permitted by law.\n", "v1.3.7")
}

func printHelpMessage() {
	fmt.Println("List information about the FILEs with ICONS and GIT STATUS (the current dir \nby default). Sort entries alphabetically if none of -tvSUX is specified.")

	getopt.PrintUsage(os.Stdout)

	fmt.Println("\nPossible value for --time-style (-T)")
	fmt.Printf(" %-11s %-32q\n", "ANSIC", "Mon Jan _2 15:04:05 2006")
	fmt.Printf(" %-11s %-32q\n", "UnixDate", "Mon Jan _2 15:04:05 MST 2006")
	fmt.Printf(" %-11s %-32q\n", "RubyDate", "Mon Jan 02 15:04:05 -0700 2006")
	fmt.Printf(" %-11s %-32q\n", "RFC822", "02 Jan 06 15:04 MST")
	fmt.Printf(" %-11s %-32q\n", "RFC822Z", "02 Jan 06 15:04 -0700")
	fmt.Printf(" %-11s %-32q\n", "RFC850", "Monday, 02-Jan-06 15:04:05 MST")
	fmt.Printf(" %-11s %-32q\n", "RFC1123", "Mon, 02 Jan 2006 15:04:05 MST")
	fmt.Printf(" %-11s %-32q\n", "RFC1123Z", "Mon, 02 Jan 2006 15:04:05 -0700")
	fmt.Printf(" %-11s %-32q\n", "RFC3339", "2006-01-02T15:04:05Z07:00")
	fmt.Printf(" %-11s %-32q\n", "Kitchen", "3:04PM")
	fmt.Printf(" %-11s %-32q [Default]\n", "Stamp", "Mon Jan _2 15:04:05")
	fmt.Printf(" %-11s %-32q\n", "StampMilli", "Jan _2 15:04:05.000")

	fmt.Println("\nExit status:")
	fmt.Println(" 0  if OK,")
	fmt.Println(" 1  if minor problems (e.g., cannot access subdirectory),")
	fmt.Println(" 2  if serious trouble (e.g., cannot access command-line argument).")
}

func getCustomTerminalWidth() int {
	// screen width for custom tw
	w, _, e := term.GetSize(int(os.Stdout.Fd()))

	if e != nil {
		return standardTerminalWidth
	}

	if w == 0 {
		// for systems that don’t support ‘TIOCGWINSZ’.
		w, _ = strconv.Atoi(os.Getenv("COLUMNS"))
	}

	return w
}
