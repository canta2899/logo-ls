package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/canta2899/logo-ls/ctw"
	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/model"
)

// App represents the main application that holds configuration, a writer, exit codes, and a logger.
type App struct {
	Config   *Config
	Writer   io.Writer
	ExitCode model.ExitCode
	Logger   *log.Logger
}

// Args stores the parsed command-line arguments as separate files and directories.
type Args struct {
	Files []model.FileEntry
	Dirs  []model.DirectoryEntry
}

// Exit terminates the program with the current stored exit code.
func (a *App) Exit() {
	os.Exit(int(a.ExitCode))
}

// Write copies a buffer to the App's writer, panicking on error.
func (a *App) Write(buf *bytes.Buffer) {
	if _, err := io.Copy(a.Writer, buf); err != nil {
		panic(err)
	}
}

// GetArguments parses the configured file list and categorizes them as files or directories.
// It also sets an error exit code for entries that are not accessible.
func (a *App) GetArguments() *Args {
	// Sort all user inputs.
	sort.Slice(a.Config.FileList, func(i, j int) bool {
		return a.Config.FileList[i] < a.Config.FileList[j]
	})

	args := &Args{}

	for _, argPath := range a.Config.FileList {
		abs, err := filepath.Abs(argPath)
		if err != nil {
			a.Logger.Printf("cannot get absolute path for %q: %v\n", argPath, err)
			a.ExitCode.SetSerious()
			continue
		}

		f, err := os.Open(abs)
		if err != nil {
			a.Logger.Printf("cannot access %q: %v\n", argPath, err)
			a.ExitCode.SetSerious()
			continue
		}

		fi, err := f.Stat()
		if err != nil {
			a.Logger.Printf("cannot stat %q: %v\n", argPath, err)
			f.Close()
			a.ExitCode.SetSerious()
			continue
		}

		if fi.IsDir() {
			args.Dirs = append(args.Dirs, model.DirectoryEntry{
				File:    *f,
				AbsPath: abs,
			})
		} else {
			args.Files = append(args.Files, model.FileEntry{
				FileInfo: fi,
				AbsPath:  abs,
			})
			f.Close() // no need to keep file handles for single files
		}
	}
	return args
}

// Run is the main entry point that orchestrates listing files/directories and printing results.
func (a *App) Run() {
	args := a.GetArguments()

	// 1. Process and display all files first.
	if len(args.Files) > 0 {
		filesDir := a.ProcessFiles(args.Files)
		a.PrintDirectory(filesDir)
		if len(args.Dirs) > 0 {
			fmt.Println()
		}
	}

	// 2. Process and display all directories (recursively or not).
	if a.Config.Recursive {
		a.processDirsRecursively(args.Dirs)
	} else {
		a.processDirsNonRecursively(args.Dirs)
	}
}

// processDirsRecursively prints each directory and its subdirectories.
func (a *App) processDirsRecursively(dirs []model.DirectoryEntry) {
	currentAbs, _ := filepath.Abs(".")
	openDirIcon := model.OpenDirIcon
	if a.Config.DisableIcon {
		openDirIcon = ""
	}

	for i, dirEntry := range dirs {
		if i > 0 {
			fmt.Println()
		}

		relName := dirEntry.Name()
		if rel, err := filepath.Rel(currentAbs, dirEntry.Name()); err == nil {
			relName = rel
		}

		fmt.Printf("%s:\n", openDirIcon+relName)

		a.recurseDirectory(&dirEntry, currentAbs)
	}
}

// processDirsNonRecursively prints each directory (but not subdirectories).
func (a *App) processDirsNonRecursively(dirs []model.DirectoryEntry) {
	pName := len(dirs) > 1
	openDirIcon := model.OpenDirIcon
	if a.Config.DisableIcon {
		openDirIcon = ""
	}

	for i, dirEntry := range dirs {
		if pName {
			fmt.Printf("%s:\n", openDirIcon+dirEntry.Name())
		}

		d, err := a.ProcessDirectory(&dirEntry)
		dirEntry.Close()
		if err != nil {
			a.Logger.Printf("partial access to %q: %v\n", dirEntry.Name(), err)
			a.ExitCode.SetSerious()
		}

		a.PrintDirectory(d)
		if i < len(dirs)-1 {
			fmt.Println()
		}
	}
}

// RecurseDirectory is the helper that processes a directory, prints it, and
// recurses through subdirectories if any.
func (a *App) recurseDirectory(dir *model.DirectoryEntry, startingAbsolutePath string) {
	d, err := a.ProcessDirectory(dir)
	dir.Close()
	if err != nil {
		a.Logger.Printf("partial access to %q: %v\n", dir.Name(), err)
		a.ExitCode.SetMinor()
	}

	a.PrintDirectory(d)

	// Recurse into subdirectories
	if len(d.Dirs) == 0 {
		return
	}

	sort.Strings(d.Dirs)
	for _, subdir := range d.Dirs {
		// Attempt to compute relative path for printing
		childPath := filepath.Join(dir.Name(), subdir)
		if rel, err := filepath.Rel(startingAbsolutePath, childPath); err == nil {
			childPath = rel
		}

		fmt.Printf("\n%s:\n", model.OpenDirIcon+childPath)

		f, err := os.Open(filepath.Join(dir.Name(), subdir))
		if err != nil {
			a.Logger.Printf("cannot access %q: %v\n", childPath, err)
			a.ExitCode.SetMinor()
			continue
		}
		abs, err := filepath.Abs(filepath.Join(dir.Name(), subdir))
		if err != nil {
			a.Logger.Println("Cannot compute abs path for:", childPath)
			f.Close()
			continue
		}
		next := &model.DirectoryEntry{File: *f, AbsPath: abs}
		a.recurseDirectory(next, startingAbsolutePath)
	}
}

// ProcessFiles converts a slice of file entries into a *model.Directory for printing.
func (a *App) ProcessFiles(files []model.FileEntry) *model.Directory {
	t := new(model.Directory)
	isLong := a.Config.LongListingMode != model.LongListingNone

	for _, fileEntry := range files {
		entry := a.buildEntry(fileEntry.AbsPath, fileEntry.FileInfo, isLong)
		t.Files = append(t.Files, entry)
	}

	return t
}

// ProcessDirectory reads the contents of the given directory, builds a *model.Directory
// that contains *model.Entry objects for each item, and returns it.
func (a *App) ProcessDirectory(d *model.DirectoryEntry) (*model.Directory, error) {
	defer func() {
		// We defer the close in case 'Readdir' triggers partial reads.
		// We manually call d.Close() in some error paths, but a double-close is no-op for os.File.
		_ = d.Close()
	}()

	dirStat, err := d.Stat()
	if err != nil {
		return nil, err
	}

	// Populate the directory structure (all internal files).
	dirModel, err := a.populateDirectory(d, dirStat)
	return dirModel, err
}

// populateDirectory reads directory contents (unless -d is set), creates *model.Entry objects,
// adds special entries (.), (..) if needed, and sets up Git statuses if requested.
func (a *App) populateDirectory(d *model.DirectoryEntry, dirStat os.FileInfo) (*model.Directory, error) {
	t := new(model.Directory)
	isLong := a.Config.LongListingMode != model.LongListingNone

	// If we need to show the current directory as an entry
	if a.Config.AllMode == model.IncludeAll || a.Config.Directory {
		t.Info = a.buildEntry(d.Name(), dirStat, isLong)
		t.Info.Name = "." // override the name to "."

		if !a.Config.DisableIcon && !a.Config.Directory {
			t.Info.Icon = icons.IconDef["diropen"].GetGlyph()
			t.Info.IconColor = icons.IconDef["diropen"].GetColor(1)
		}
	}

	// If -d is set, we only show the directory itself, not its contents.
	if a.Config.Directory {
		t.Files = append(t.Files, t.Info)
		return t, nil
	}

	// Perform Readdir
	files, err := d.Readdir(0)
	if err != nil {
		// Return partial result + error
		return t, err
	}

	// If Git status is requested, prepare the repository info map.
	var gitRepoStatus map[string]string
	if a.Config.GitStatus {
		// might be nil if not a git repo or no statuses found
		gitRepoStatus = format.GetFilesGitStatus(d.Name())
		if len(gitRepoStatus) == 0 {
			gitRepoStatus = nil
		}
	}

	showHidden := a.Config.AllMode != model.IncludeDefault

	// Build entries for each file
	for _, fi := range files {
		name := fi.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fullpath := filepath.Join(d.Name(), name)
		entry := a.buildEntry(fullpath, fi, isLong)

		// Overwrite the name field if necessary
		// (because buildEntry sets Name/Ext from fi.Name())
		entry.Indicator = format.GetIndicator(fullpath, isLong)

		// Overwrite icons for symlinks to reflect target type
		if !a.Config.DisableIcon {
			entry.Icon, entry.IconColor = format.GetIcon(name, entry.Ext, entry.Indicator)
			if format.IsLink(fullpath) {
				if s, err := filepath.EvalSymlinks(fullpath); err == nil {
					realExt := filepath.Ext(s)
					realName := s[0 : len(s)-len(realExt)]
					realIndicator := format.GetIndicator(s, isLong)
					entry.Icon, entry.IconColor = format.GetIcon(realName, realExt, realIndicator)
					if a.Config.DisableColor {
						entry.IconColor = ""
					}
				}
			}
		}

		// If Git status is available, attach it.
		if gitRepoStatus != nil {
			// Directories in Git can be recognized by trailing slash in the map.
			if fi.IsDir() {
				fmt.Println(fi.Name() + model.PathSeparator)
				entry.GitStatus = gitRepoStatus[fi.Name()+model.PathSeparator]
			} else {
				entry.GitStatus = gitRepoStatus[fi.Name()]
			}
		}

		t.Files = append(t.Files, entry)
		if fi.IsDir() {
			t.Dirs = append(t.Dirs, name+"/")
		}
	}

	// If -a (IncludeAll) is passed, add "." and ".." entries
	if a.Config.AllMode == model.IncludeAll {
		// t.Info is ".", so we add that to the Files if not already present
		if t.Info != nil {
			t.Files = append(t.Files, t.Info)
		}

		// Build ".." parent directory entry
		pp := filepath.Dir(d.Name())
		pStat, err2 := os.Lstat(pp)
		if err2 == nil { // if we can't stat parent, skip adding
			parentEntry := a.buildEntry(pp, pStat, isLong)
			parentEntry.Name = ".."
			// Overwrite icon for parent
			if !a.Config.DisableIcon {
				parentEntry.Icon = icons.IconDef["diropen"].GetGlyph()
				if !a.Config.DisableColor {
					parentEntry.IconColor = icons.IconDef["diropen"].GetColor(1)
				}
			}
			t.Files = append(t.Files, parentEntry)
			t.Parent = parentEntry
		}
	}
	return t, err
}

// buildEntry is a helper that constructs a *model.Entry from a given path, os.FileInfo, and
// whether we are in a long-listing context.
func (a *App) buildEntry(fullPath string, fi os.FileInfo, isLong bool) *model.Entry {
	entry := &model.Entry{}

	// Basic name, extension, size, mod time
	name := fi.Name()
	entry.Ext = filepath.Ext(name)
	entry.Name = name[0 : len(name)-len(entry.Ext)]
	entry.Size = fi.Size()
	entry.ModTime = fi.ModTime()

	// Inode
	if a.Config.ShowInodeNumber {
		entry.InodeNumber = format.GetInodeNumber(fullPath)
	}

	// File mode/permissions
	if isLong {
		entry.Mode = fi.Mode().String()
		entry.ModeBits = uint32(fi.Mode())
		entry.NumHardLinks = format.GetHardLinkCount(fullPath)
		entry.Owner, entry.Group = model.GetOwnerGroupInfo(fi, a.Config.NoGroup, a.Config.LongListingMode)
	}

	// Block size
	if a.Config.ShowBlockSize {
		model.DirBlocks(entry, fi)
	}

	// **FIX for directory icons**.
	if !a.Config.DisableIcon {
		entry.Icon, entry.IconColor = format.GetIcon(entry.Name, entry.Ext, entry.Indicator)

		if a.Config.DisableColor {
			entry.IconColor = ""
		}
	}

	return entry
}

// PrintDirectory sorts the directory's files according to the app config and prints them.
func (a *App) PrintDirectory(d *model.Directory) {
	if d == nil {
		return
	}

	// Sort the directory contents
	format.SetLessFunction(d, a.Config.SortMode)
	d.Sort(a.Config.SortMode, a.Config.Reverse)

	// Prepare the columnar text writer
	lineCtw := a.getCTW()

	// Decide printing format (long, single line, or multi-column)
	switch {
	case a.Config.LongListingMode != model.LongListingNone:
		for _, f := range d.Files {
			lineCtw.AddRow(
				f.IconColor,
				a.blockSizeWithInode(f),
				f.Mode,
				strconv.Itoa(int(f.NumHardLinks)),
				f.Owner,
				f.Group,
				format.GetFormattedSize(f.Size, a.Config.HumanReadable),
				a.Config.TimeFormatter.Format(&f.ModTime),
				f.Icon,
				f.Name+f.Ext+f.Indicator,
				f.GitStatus,
			)
		}

	case a.Config.OneFilePerLine:
		for _, f := range d.Files {
			lineCtw.AddRow(
				f.IconColor,
				a.blockSizeWithInode(f),
				f.Icon,
				f.Name+f.Ext+f.Indicator,
				f.GitStatus,
			)
		}

	default:
		for _, f := range d.Files {
			lineCtw.AddRow(
				f.IconColor,
				a.blockSizeWithInode(f),
				f.Icon,
				f.Name+f.Ext+f.Indicator,
				f.GitStatus,
			)
		}
	}

	// Flush to buffer and write out
	buf := new(bytes.Buffer)
	lineCtw.Flush(buf)
	a.Write(buf)
}

// getCTW instantiates the proper ctw.CTW based on app config.
func (a *App) getCTW() ctw.CTW {
	var out ctw.CTW

	switch {
	case a.Config.LongListingMode != model.LongListingNone:
		out = ctw.NewLongCTW(10)
	case a.Config.OneFilePerLine:
		out = ctw.NewLongCTW(4)
	default:
		out = ctw.NewStandardCTW(a.Config.TerminalWidth)
	}

	if a.Config.DisableColor {
		out.DisplayColor(false)
		model.OpenDirIcon = icons.IconDef["diropen"].GetGlyph() + " "
	}

	if a.Config.DisableIcon {
		model.OpenDirIcon = ""
	}

	return out
}

// blockSizeWithInode generates the block size and optional inode as a single string.
func (a *App) blockSizeWithInode(e *model.Entry) string {
	var parts []string

	if a.Config.ShowInodeNumber {
		parts = append(parts, e.InodeNumber)
	}

	if a.Config.ShowBlockSize {
		parts = append(parts, format.GetFormattedSize(e.Blocks, a.Config.HumanReadable))
	}

	return strings.TrimSpace(strings.Join(parts, " "))
}
